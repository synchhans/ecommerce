package user

import "context"

type Repository interface {
	CreateUser(ctx context.Context, email, passwordHash, name string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, string /*passwordHash*/, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
}
