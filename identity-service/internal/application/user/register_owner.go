package user

import (
	"context"
	"encoding/json"
	"identity-service/internal/domain/auth"
	"identity-service/internal/domain/outbox"
	"identity-service/internal/domain/user"
	"identity-service/internal/shared/event"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RegisterOwner struct {
	db            *gorm.DB
	emailVerifier auth.EmailVerifier
	hasher        auth.PasswordHasher
	repo          user.UserRepository
	outboxRepo    outbox.OutboxRepository
	publisher     event.EventPublisher
}

func NewRegisterOwner(
	db *gorm.DB,
	emailVerifier auth.EmailVerifier,
	hasher auth.PasswordHasher,
	repo user.UserRepository,
	outboxRepo outbox.OutboxRepository,
	publisher event.EventPublisher,
) *RegisterOwner {
	return &RegisterOwner{
		db:            db,
		emailVerifier: emailVerifier,
		hasher:        hasher,
		repo:          repo,
		outboxRepo:    outboxRepo,
		publisher:     publisher,
	}
}

func (uc *RegisterOwner) Execute(ctx context.Context, input RegisterOwnerRequest) (Response, error) {
	if err := uc.emailVerifier.Verify(ctx, input.Email, input.Code); err != nil {
		return Response{}, err
	}

	hashedPassword, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return Response{}, err
	}

	userID, err := uuid.NewV7()
	if err != nil {
		return Response{}, err
	}

	restaurantID, err := uuid.NewV7()
	if err != nil {
		return Response{}, err
	}

	newUser := user.User{
		ID:        userID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashedPassword,
		Role:      input.Role,
		Status:    user.DefaultStatus,
	}

	payloadMap := map[string]interface{}{
		"restaurant_id": restaurantID,
		"owner_id":      userID,
		"business_name": input.BusinessName,
		"vat_number":    input.VATNumber,
	}

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return Response{}, err
	}

	err = uc.db.Transaction(func(tx *gorm.DB) error {
		if err := uc.repo.WithTx(tx).Create(ctx, &newUser); err != nil {
			return err
		}

		newEvent := outbox.NewOutboxEvent(
			restaurantID,
			outbox.EventRestaurantInitiated,
			payload,
		)

		if err := uc.outboxRepo.WithTx(tx).Create(ctx, &newEvent); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return Response{}, err
	}

	t := time.Now().UTC()
	userRegistered := UserRegistered{
		Email:     newUser.Email,
		FirstName: newUser.FirstName,
		Role:      newUser.Role,
		Timestamp: t.Format(time.RFC3339),
	}
	userRegistered.EventName = userRegistered.GetEventName()

	if err := uc.publisher.Publish(userRegistered); err != nil {
		log.Println("Failed to publish user.registered event:", err)
	}

	return MapToResponse(&newUser), nil
}
