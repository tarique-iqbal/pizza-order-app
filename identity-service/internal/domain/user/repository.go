package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	EmailExists(email string) (bool, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
}
