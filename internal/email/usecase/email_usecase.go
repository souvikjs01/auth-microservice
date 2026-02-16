package usecase

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/souvikjs01/auth-microservice/internal/email"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	"github.com/streadway/amqp"
)

type EmailUsecase struct {
	emailRepo email.EmailRepository
	logger    logger.Logger
}

func NewEmailUsecase(emailRepo email.EmailRepository, logger logger.Logger) *EmailUsecase {
	return &EmailUsecase{
		emailRepo: emailRepo,
		logger:    logger,
	}
}

func (e *EmailUsecase) SendEmail(ctx context.Context, delivery amqp.Delivery) error {
	reader := bytes.NewReader(delivery.Body)
	deliveryBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadAll")

	}

	e.logger.Info("Sending email to user with body: %s", string(deliveryBytes))
	return nil
}
