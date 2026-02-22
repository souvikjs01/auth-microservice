package models

import (
	"time"

	"github.com/google/uuid"
)

type Email struct {
	EmailID     uuid.UUID `json:"email_id" db:"email_id" validate:"omitempty"`
	To          []string  `json:"to" db:"to" validate:"required"`
	From        string    `json:"from" db:"from" validate:"required,email"`
	Body        string    `json:"body" db:"body" validate:"required"`
	Subject     string    `json:"subject" db:"subject" validate:"required"`
	ContentType string    `json:"content_type" db:"content_type" validate:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" validate:"omitempty"`
}
