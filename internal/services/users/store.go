package users

import (
	"context"

	"github.com/aschles4/finalProject/internal/pkg/dynamo/users"
)

type store interface {
	CreateUserAccount(ctx context.Context, account users.UserAccount) error
	CreateUserProfile(ctx context.Context, profile users.UserProfile) error
	FindUserAccountByEmail(ctx context.Context, email string) (*users.UserAccount, error)
	FindUserAccountByToken(ctx context.Context, token string) (*users.UserAccount, error)
	FindUserProfileByID(ctx context.Context, ID string) (*users.UserProfile, error)
	RemoveUserAccountByID(ctx context.Context, ID string) error
	RemoveUserProfileByID(ctx context.Context, ID string) error
	UpdateUserAccountToken(ctx context.Context, ID, token string) (*users.UserAccount, error)
}
