package user

type UserCreatedEvent struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Timestamp string `json:"timestamp"`
}

func (e UserCreatedEvent) GetEventName() string {
	return "user.created"
}
