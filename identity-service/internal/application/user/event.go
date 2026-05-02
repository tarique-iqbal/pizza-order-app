package user

import "time"

type UserRegistered struct {
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	Role      string    `json:"role"`
	EventName string    `json:"event_name"`
	Timestamp time.Time `json:"timestamp"`
}

func (e UserRegistered) GetEventName() string {
	return "user.registered"
}
