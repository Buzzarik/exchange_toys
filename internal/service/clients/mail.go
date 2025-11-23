package clients

import (
	"service/internal/models"
	"service/internal/service"

	"net/smtp"
	"fmt"
	"log/slog"
	"time"
)

func SendEmailToSingleParticipant(app *service.Application, srcUserId string, dstUserId string) {
	const maxRetries = 3
	const retryDelay = 2 * time.Second

	// Получаем данные пользователя
	dbSrcUser, err := app.Storage.SelectUserById(&models.User{UserId: srcUserId})
	if err != nil || dbSrcUser == nil {
		app.Log.Error("Failed to get user email", 
			slog.String("user_id", srcUserId), 
			slog.Any("error", err))
		return
	}

	dbDstUser, err := app.Storage.SelectUserById(&models.User{UserId: dstUserId})
	if err != nil || dbDstUser == nil {
		app.Log.Error("Failed to get user email", 
			slog.String("user_id", dstUserId), 
			slog.Any("error", err))
		return
	}

	// Пытаемся отправить письмо с повторными попытками
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := sendSingleEmail(dbSrcUser, dbDstUser)
		if err == nil {
			app.Log.Info("Email sent successfully to participant",
				slog.String("user_id", dbSrcUser.UserId))
			return
		}

		app.Log.Warn("Failed to send email to participant, retrying",
			slog.String("user_id", dbSrcUser.UserId),
			slog.Int("attempt", attempt),
			slog.Any("error", err))

		if attempt < maxRetries {
			time.Sleep(retryDelay * time.Duration(attempt))
		}
	}

	app.Log.Error("Failed to send email to participant after all retries",
		slog.String("user_id", dbSrcUser.UserId))
}

func sendSingleEmail(srcUser *models.User, dstUser *models.User) error {
	subject := fmt.Sprintf("С вами хотят связаться",
	)
	body := fmt.Sprintf(
		"Уважаемый(ая) %s!\n\nВы можете списаться и обсудить детали обмена с %s по почте %s",
		srcUser.FullName(),
		dstUser.FullName(),
		dstUser.Email,
	)

	// Используем ваш MailHog или SMTP
	return sendViaMailHog(srcUser.Email, subject, body)
}

func sendViaMailHog(to, subject, body string) error {
    const from = "exchangeToy@yandex.ru"
	const host = "localhost:1025"
    
    message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", 
        from, to, subject, body)

    return smtp.SendMail(
        host, // MailHog SMTP порт
        nil,              // аутентификация не требуется
        from,
        []string{to},
        []byte(message),
    )
}