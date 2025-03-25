package user_test

import (
	"errors"
	"testing"

	aUser "pizza-order-api/internal/application/user"
	"pizza-order-api/internal/domain/user"
	"pizza-order-api/internal/infrastructure/security"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestSignInUserUseCase_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	password, _ := security.HashPassword("password")

	mockRepo.On("FindByEmail", "test@example.com").Return(&user.User{
		ID:       1,
		Email:    "test@example.com",
		Password: password,
	}, nil)

	useCase := aUser.NewSignInUserUseCase(mockRepo)

	token, err := useCase.Execute("test@example.com", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestSignInUserUseCase_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	password, _ := security.HashPassword("password")

	mockRepo.On("FindByEmail", "test@example.com").Return(&user.User{
		ID:       1,
		Email:    "test@example.com",
		Password: password,
	}, nil)

	useCase := aUser.NewSignInUserUseCase(mockRepo)

	token, err := useCase.Execute("test@example.com", "wrongpassword")
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestSignInUserUseCase_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, errors.New("user not found"))

	useCase := aUser.NewSignInUserUseCase(mockRepo)

	token, err := useCase.Execute("notfound@example.com", "password")
	assert.Error(t, err)
	assert.Empty(t, token)
}
