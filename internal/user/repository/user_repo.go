package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/souvikjs01/auth-microservice/internal/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.Create")
	defer span.Finish()

	query := `INSERT INTO users (first_name, last_name, email, password, role, avatar) VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, ''), 'user')::role, $6) RETURNING *`

	createdUser := &models.User{}

	if err := u.db.QueryRowxContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password, user.Role, user.Avatar).StructScan(createdUser); err != nil {
		return nil, errors.Wrap(err, "Register.QueryRowxContext")
	}

	return createdUser, nil
}

// find by email
func (u *UserRepository) FindBYEmail(ctx context.Context, email string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.FindBYEmail")
	defer span.Finish()

	query := `SELECT user_id, first_name, last_name, email, role, avatar, password, created_at, updated_at
	          FROM users 
			  WHERE email=$1`

	user := &models.User{}

	if err := u.db.GetContext(ctx, user, query, email); err != nil {
		return nil, errors.Wrap(err, "FindBYEmail.GetContext")
	}
	return user, nil
}

// find by id
func (u *UserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.FindById")
	defer span.Finish()

	query := `SELECT user_id, first_name, last_name, email, role, avatar, created_at, updated_at
			  FROM users
			  WHERE user_id=$1`

	user := &models.User{}

	if err := u.db.GetContext(ctx, user, query, userID); err != nil {
		return nil, errors.Wrap(err, "FindByID.GetContext")
	}

	return user, nil
}
