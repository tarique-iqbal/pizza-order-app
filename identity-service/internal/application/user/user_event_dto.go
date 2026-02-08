package user

type UserCreatedEvent struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Role      string `json:"role"`
	EventName string `json:"event_name"`
	Timestamp string `json:"timestamp"`
}

func (e UserCreatedEvent) GetEventName() string {
	return "user.registered"
}
