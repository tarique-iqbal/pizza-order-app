package user

import (
	"context"
	"errors"

	"identity-service/internal/domain/user"

	"github.com/google/uuid"
)

type FindByID struct {
	repo user.UserRepository
}

func NewFindByID(repo user.UserRepository) *FindByID {
	return &FindByID{repo: repo}
}

func (uc *FindByID) Execute(ctx context.Context, userID uuid.UUID) (Response, error) {
	user, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return Response{}, errors.New("internal server error")
	}

	if user == nil {
		return Response{}, errors.New("user not found")
	}

	return MapToResponse(user), nil
}
