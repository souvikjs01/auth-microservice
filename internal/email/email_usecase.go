package email

import (
	"context"

	"github.com/streadway/amqp"
)

type EmailUsecase interface {
	SendEmail(ctx context.Context, delivery amqp.Delivery) error
}
