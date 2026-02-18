package email

import (
	"context"

	"github.com/souvikjs01/auth-microservice/internal/models"
)

type Mailer interface {
	Send(ctx context.Context, email *models.Email) error
}
