package mailer

import (
	"github.com/souvikjs01/auth-microservice/config"
	"gopkg.in/gomail.v2"
)

func NewMailDialer(cfg *config.Config) *gomail.Dialer {
	return gomail.NewDialer(
		cfg.Smtp.Host,
		cfg.Smtp.Port,
		cfg.Smtp.User,
		cfg.Smtp.Password,
	)
}
