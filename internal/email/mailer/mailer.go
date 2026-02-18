package mailer

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/internal/models"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	cfg        *config.Config
	mailDialer *gomail.Dialer
}

func NewMailer(cfg *config.Config, mailDialer *gomail.Dialer) *Mailer {
	return &Mailer{
		cfg:        cfg,
		mailDialer: mailDialer,
	}
}

func (m *Mailer) Send(ctx context.Context, email *models.Email) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Mailer.Send")
	defer span.Finish()

	msg := gomail.NewMessage()
	msg.SetHeader("From", email.From)
	msg.SetHeader("To", email.To...)
	msg.SetHeader("Subject", email.Subject)
	msg.SetBody(email.ContentType, email.Body)

	return m.mailDialer.DialAndSend(msg)
}
