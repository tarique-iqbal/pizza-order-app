package user

import (
	"context"
	"errors"

	domain "identity-service/internal/domain/user"
)

type FindByID struct {
	repo domain.UserRepository
}

func NewFindByID(repo domain.UserRepository) *FindByID {
	return &FindByID{repo: repo}
}

func (uc *FindByID) Execute(ctx context.Context, id int) (Response, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return Response{}, errors.New("internal server error")
	}

	if user == nil {
		return Response{}, errors.New("user not found")
	}

	return MapToResponse(user), nil
}
