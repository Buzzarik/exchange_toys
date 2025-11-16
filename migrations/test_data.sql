-- Очистка таблиц
TRUNCATE TABLE exchange_details CASCADE;
TRUNCATE TABLE exchange CASCADE;
TRUNCATE TABLE toys CASCADE;
TRUNCATE TABLE users CASCADE;

-- Вставляем тестовых пользователей
INSERT INTO users (user_id, first_name, middle_name, last_name, email, password_hash) VALUES
('user_1', 'Иван', 'Иванович', 'Иванов', 'ivan@example.com', 'hash_user_1'),
('user_2', 'Петр', 'Петрович', 'Петров', 'petr@example.com', 'hash_user_2'),
('user_3', 'Мария', 'Сергеевна', 'Сидорова', 'maria@example.com', 'hash_user_3'),
('user_4', 'Анна', 'Владимировна', 'Кузнецова', 'anna@example.com', 'hash_user_4')
ON CONFLICT (user_id) DO NOTHING;

-- Вставляем тестовые игрушки
INSERT INTO toys (toy_id, user_id, name, description, photo_url, idempotency_token, status) VALUES
-- Пользователь 1: несколько игрушек
('toy_1', 'user_1', 'Машинка', 'Красная машинка', '/photos/car.jpg', 'token_toy_1', 'created'),
('toy_2', 'user_1', 'Кукла', 'Барби', '/photos/doll.jpg', 'token_toy_2', 'created'),
('toy_3', 'user_1', 'Конструктор', 'Лего', '/photos/lego.jpg', 'token_toy_3', 'created'),

-- Пользователь 2: несколько игрушек  
('toy_4', 'user_2', 'Мяч', 'Футбольный мяч', '/photos/ball.jpg', 'token_toy_4', 'created'),
('toy_5', 'user_2', 'Пазл', 'Детский пазл', '/photos/puzzle.jpg', 'token_toy_5', 'created'),

-- Пользователь 3: несколько игрушек
('toy_6', 'user_3', 'Кубики', 'Деревянные кубики', '/photos/blocks.jpg', 'token_toy_6', 'created'),
('toy_7', 'user_3', 'Мишка', 'Плюшевый мишка', '/photos/bear.jpg', 'token_toy_7', 'created'),
('toy_8', 'user_3', 'Раскраска', 'Детская раскраска', '/photos/coloring.jpg', 'token_toy_8', 'created'),

-- Пользователь 4: несколько игрушек
('toy_9', 'user_4', 'Матрёшка', 'Русская матрёшка', '/photos/matryoshka.jpg', 'token_toy_9', 'created'),
('toy_10', 'user_4', 'Юла', 'Игрушка-юла', '/photos/spinner.jpg', 'token_toy_10', 'created')
ON CONFLICT (toy_id) DO NOTHING;

-- 1) Несколько обменов с одной игрушкой (toy_1 участвует в нескольких обменах)
INSERT INTO exchange (exchange_id, src_toy_id, dst_toy_id, idempotency_token, status) VALUES
('exchange_1', 'toy_1', 'toy_4', 'token_exchange_1', 'created'), -- user_1 отдает toy_1 за toy_4 user_2
('exchange_2', 'toy_1', 'toy_6', 'token_exchange_2', 'created'), -- user_1 отдает toy_1 за toy_6 user_3
('exchange_3', 'toy_1', 'toy_9', 'token_exchange_3', 'created')  -- user_1 отдает toy_1 за toy_9 user_4
ON CONFLICT (exchange_id) DO NOTHING;

-- Детали обменов для scenario 1
INSERT INTO exchange_details (exchange_id, toy_id, user_id, status) VALUES
-- Обмен 1: toy_1 ↔ toy_4
('exchange_1', 'toy_1', 'user_1', 'created'),
('exchange_1', 'toy_4', 'user_2', 'created'),

-- Обмен 2: toy_1 ↔ toy_6  
('exchange_2', 'toy_1', 'user_1', 'created'),
('exchange_2', 'toy_6', 'user_3', 'created'),

-- Обмен 3: toy_1 ↔ toy_9
('exchange_3', 'toy_1', 'user_1', 'created'),
('exchange_3', 'toy_9', 'user_4', 'created')
ON CONFLICT (exchange_id, user_id, toy_id) DO NOTHING;

-- 2) Несколько обменов с разными игрушками одного пользователя (user_1 участвует в разных обменах разными игрушками)
INSERT INTO exchange (exchange_id, src_toy_id, dst_toy_id, idempotency_token, status) VALUES
('exchange_4', 'toy_2', 'toy_5', 'token_exchange_4', 'created'), -- user_1 отдает toy_2 за toy_5 user_2
('exchange_5', 'toy_3', 'toy_7', 'token_exchange_5', 'created'), -- user_1 отдает toy_3 за toy_7 user_3
('exchange_6', 'toy_2', 'toy_10', 'token_exchange_6', 'created') -- user_1 отдает toy_2 за toy_10 user_4
ON CONFLICT (exchange_id) DO NOTHING;

-- Детали обменов для scenario 2
INSERT INTO exchange_details (exchange_id, toy_id, user_id, status) VALUES
-- Обмен 4: toy_2 ↔ toy_5
('exchange_4', 'toy_2', 'user_1', 'created'),
('exchange_4', 'toy_5', 'user_2', 'created'),

-- Обмен 5: toy_3 ↔ toy_7
('exchange_5', 'toy_3', 'user_1', 'created'),
('exchange_5', 'toy_7', 'user_3', 'created'),

-- Обмен 6: toy_2 ↔ toy_10
('exchange_6', 'toy_2', 'user_1', 'created'),
('exchange_6', 'toy_10', 'user_4', 'created')
ON CONFLICT (exchange_id, user_id, toy_id) DO NOTHING;

-- 3) Дополнительные обмены между разными пользователями
INSERT INTO exchange (exchange_id, src_toy_id, dst_toy_id, idempotency_token, status) VALUES
('exchange_7', 'toy_8', 'toy_4', 'token_exchange_7', 'created'), -- user_3 отдает toy_8 за toy_4 user_2
('exchange_8', 'toy_9', 'toy_6', 'token_exchange_8', 'created')  -- user_4 отдает toy_9 за toy_6 user_3
ON CONFLICT (exchange_id) DO NOTHING;

-- Детали обменов для scenario 3
INSERT INTO exchange_details (exchange_id, toy_id, user_id, status) VALUES
-- Обмен 7: toy_8 ↔ toy_4
('exchange_7', 'toy_8', 'user_3', 'created'),
('exchange_7', 'toy_4', 'user_2', 'created'),

-- Обмен 8: toy_9 ↔ toy_6
('exchange_8', 'toy_9', 'user_4', 'created'),
('exchange_8', 'toy_6', 'user_3', 'created')
ON CONFLICT (exchange_id, user_id, toy_id) DO NOTHING;


SELECT * FROM toys ORDER BY user_id;

SELECT 
	e.exchange_id, 
	e.idempotency_token, 
	e.status AS exchange_status, 
	ed.toy_id, 
	ed.user_id, 
	ed.status AS user_status
FROM exchange AS e
INNER JOIN exchange_details AS ed
ON (e.exchange_id = ed.exchange_id)
ORDER BY e.exchange_id;
-- ORDER BY ed.toy_id;

UPDATE toys
SET status = 'removed'
WHERE toy_id = 'toy_1';

UPDATE exchange_details
SET status = 'failed'
WHERE user_id = 'user_2' AND exchange_id = 'exchange_1';
	