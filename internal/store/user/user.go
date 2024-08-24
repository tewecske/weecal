package user

import "context"

type User struct {
	ID       uint
	Email    string
	Password string
}

type UserStore interface {
	CreateUser(email string, password string) error
	GetUser(email string) (*User, error)
}

type key string

var userKey key = "user"

func NewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}
