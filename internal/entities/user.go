package entities

type User struct {
	ID    string
	Email string
}

func NewUserEntity(id, email string) (*User, error) {
	return &User{
		ID:    id,
		Email: email,
	}, nil
}
