package testutil

import (
	"time"

	"github.com/google/uuid"
)

func StringPtr(s string) *string {
	return &s
}

func Int16Ptr(i int16) *int16 {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func MustNewID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}
