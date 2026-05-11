package testutil

import "github.com/google/uuid"

func StringPtr(v string) *string {
	return &v
}

func Int16Ptr(v int16) *int16 {
	return &v
}

func MustNewID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}
