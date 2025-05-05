package user

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	EmailExists(email string) (bool, error)
}
