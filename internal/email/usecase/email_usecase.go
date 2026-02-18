package usecase

import (
	"context"
	"encoding/json"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/souvikjs01/auth-microservice/internal/email"
	"github.com/souvikjs01/auth-microservice/internal/models"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	"github.com/souvikjs01/auth-microservice/pkg/utils"
	"github.com/streadway/amqp"
)

type EmailUsecase struct {
	emailRepo email.EmailRepository
	logger    logger.Logger
	mailer    email.Mailer
}

func NewEmailUsecase(emailRepo email.EmailRepository, logger logger.Logger, mailer email.Mailer) *EmailUsecase {
	return &EmailUsecase{
		emailRepo: emailRepo,
		logger:    logger,
		mailer:    mailer,
	}
}

func (e *EmailUsecase) SendEmail(ctx context.Context, delivery amqp.Delivery) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailUsecase.SendEmail")
	defer span.Finish()

	mail := &models.Email{}
	if err := json.Unmarshal(delivery.Body, mail); err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}

	if err := utils.ValidateStruct(ctx, mail); err != nil {
		return errors.Wrap(err, "EmailUsecase.SendEmail.utils.ValidStruct")
	}

	e.logger.Infof("Sending email %+v", mail)

	return nil
}
