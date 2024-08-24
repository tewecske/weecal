package user

type TestUserStore struct {
	User  *User
	Error error
}

func (u *TestUserStore) CreateUser(email string, password string) error {
	return nil
}
func (u *TestUserStore) GetUser(email string) (*User, error) {
	if u.Error != nil {
		return nil, u.Error
	} else {
		return u.User, nil
	}
}
