-- Create Statuses
DO $$ BEGIN
    CREATE TYPE ToyStatus AS ENUM ('created', 'exchanging', 'removed', 'exchanged');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
    CREATE TYPE ExchangeStatus AS ENUM ('created', 'confirm', 'success', 'failed');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
    CREATE TYPE ExchangeDetailsStatus AS ENUM ('created', 'failed', 'confirm_1', 'confirm_2', 'success');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- Create Tables
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    first_name TEXT NOT NULL,
    middle_name TEXT,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS toys (
    toy_id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    photo_url TEXT,
    idempotency_token TEXT UNIQUE,
    status ToyStatus NOT NULL DEFAULT 'created',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS exchange (
    exchange_id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    src_toy_id TEXT NOT NULL,
    dst_toy_id TEXT NOT NULL,
    status ExchangeStatus NOT NULL DEFAULT 'created',
    idempotency_token TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS exchange_details (
    exchange_id TEXT NOT NULL,
    toy_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    status ExchangeDetailsStatus NOT NULL DEFAULT 'created',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (exchange_id, user_id, toy_id)
);

-- Create Trigger Functions
-- 1. toys.status → removed → все exchange_details по игрушке (не success/failed) → failed
CREATE OR REPLACE FUNCTION toys_removed_set_exchanges_failed()
RETURNS trigger AS $$
BEGIN
    IF NEW.status = 'removed' AND (OLD.status IS DISTINCT FROM NEW.status) THEN
        UPDATE exchange_details
        SET status = 'failed', updated_at = NOW()
        WHERE toy_id = NEW.toy_id
          AND status NOT IN ('failed', 'success');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. exchange_details.status → failed → остальные details этого обмена (не fail/success) → failed; exchange → failed
CREATE OR REPLACE FUNCTION detail_failed_propagate()
RETURNS trigger AS $$
BEGIN
    IF NEW.status = 'failed' AND (OLD.status IS DISTINCT FROM NEW.status) THEN
        UPDATE exchange_details
        SET status = 'failed', updated_at = NOW()
        WHERE exchange_id = NEW.exchange_id
          AND NOT (toy_id = NEW.toy_id AND user_id = NEW.user_id)
          AND status NOT IN ('failed', 'success');

        UPDATE exchange
        SET status = 'failed', updated_at = NOW()
        WHERE exchange_id = NEW.exchange_id AND status NOT IN ('failed', 'success');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 3. exchange_details.status → confirm_1 → если другая сторона тоже confirm_1, то exchange → confirm
CREATE OR REPLACE FUNCTION detail_confirm_1_update_exchange()
RETURNS trigger AS $$
DECLARE
    other_cnt INTEGER;
BEGIN
    IF NEW.status = 'confirm_1' AND (OLD.status IS DISTINCT FROM NEW.status) THEN
        SELECT count(*) INTO other_cnt
        FROM exchange_details
        WHERE exchange_id = NEW.exchange_id
          AND user_id <> NEW.user_id
          AND status = 'confirm_1';

        IF other_cnt > 0 THEN
            UPDATE exchange
            SET status = 'confirm', updated_at = NOW()
            WHERE exchange_id = NEW.exchange_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 4. exchange_details.status → confirm_2 → если другая сторона тоже confirm_2, то exchange → success
CREATE OR REPLACE FUNCTION detail_confirm_2_update_success()
RETURNS trigger AS $$
DECLARE
    other_cnt INTEGER;
BEGIN
    IF NEW.status = 'confirm_2' AND (OLD.status IS DISTINCT FROM NEW.status) THEN
        SELECT count(*) INTO other_cnt
        FROM exchange_details
        WHERE exchange_id = NEW.exchange_id
          AND user_id <> NEW.user_id
          AND status = 'confirm_2';


        IF other_cnt > 0 THEN
            UPDATE exchange
            SET status = 'success', updated_at = NOW()
            WHERE exchange_id = NEW.exchange_id;

            UPDATE exchange_details
            SET status = 'success', updated_at = NOW()
            WHERE exchange_id = NEW.exchange_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 5. exchange.status → success → swap владельцев игрушек, другие сделки с ними → failed (если не уже failed/success)
CREATE OR REPLACE FUNCTION exchange_success_swap_owners()
RETURNS trigger AS $$
BEGIN
    IF NEW.status = 'success' AND (OLD.status IS DISTINCT FROM NEW.status) THEN
        -- Создаем копию src_toy для dst_user
        INSERT INTO toys (user_id, name, description, photo_url, idempotency_token)
        SELECT 
            (SELECT user_id FROM toys WHERE toy_id = NEW.dst_toy_id),
            name, description, photo_url, gen_random_uuid()::text
        FROM toys WHERE toy_id = NEW.src_toy_id;
        
        -- Создаем копию dst_toy для src_user
        INSERT INTO toys (user_id, name, description, photo_url, idempotency_token)
        SELECT 
            (SELECT user_id FROM toys WHERE toy_id = NEW.src_toy_id),
            name, description, photo_url, gen_random_uuid()::text
        FROM toys WHERE toy_id = NEW.dst_toy_id;
        
        -- Помечаем оригинальные игрушки как removed
        UPDATE toys SET status = 'exchanged'
        WHERE toy_id IN (NEW.src_toy_id, NEW.dst_toy_id);
        
        -- Отменяем другие сделки с оригинальными игрушками
        UPDATE exchange_details SET status = 'failed'
        WHERE (toy_id = NEW.src_toy_id OR toy_id = NEW.dst_toy_id)
            AND exchange_id != NEW.exchange_id
            AND status NOT IN ('failed', 'success');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 6. exchange_details.status → уже success/failed → не обновляем
CREATE OR REPLACE FUNCTION prevent_update_completed_exchange_details()
RETURNS trigger AS $$
BEGIN
    -- Если пытаемся изменить запись со статусом 'failed' или 'success'
    IF OLD.status IN ('failed', 'success') THEN
        -- Ничего не делаем, возвращаем старое значение
        RETURN OLD;
    END IF;
    
    -- Иначе разрешаем обновление
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create Triggers
-- 1
CREATE TRIGGER tg_toys_removed
AFTER UPDATE OF status ON toys
FOR EACH ROW
EXECUTE FUNCTION toys_removed_set_exchanges_failed();

-- 2
CREATE TRIGGER tg_details_failed
AFTER UPDATE OF status ON exchange_details
FOR EACH ROW
EXECUTE FUNCTION detail_failed_propagate();

-- 3
CREATE TRIGGER tg_details_confirm1
AFTER UPDATE OF status ON exchange_details
FOR EACH ROW
EXECUTE FUNCTION detail_confirm_1_update_exchange();

-- 4
CREATE TRIGGER tg_details_confirm2
AFTER UPDATE OF status ON exchange_details
FOR EACH ROW
EXECUTE FUNCTION detail_confirm_2_update_success();

-- 5
CREATE TRIGGER tg_exchange_success
AFTER UPDATE OF status ON exchange
FOR EACH ROW
EXECUTE FUNCTION exchange_success_swap_owners();

-- 6
CREATE TRIGGER prevent_update_completed_exchange_details_trigger
    BEFORE UPDATE ON exchange_details
    FOR EACH ROW
    EXECUTE FUNCTION prevent_update_completed_exchange_details();