package users

import (
	"context"

	"github.com/aschles4/finalProject/internal/pkg/dynamo/users"
)

func (u Users) CreateUserProfile(ctx context.Context, ID, name, email string, acts []users.StreamAccounts) error {
	p := users.UserProfile{
		ID:             ID,
		Name:           name,
		Email:          email,
		StreamAccounts: acts,
		Library: users.Library{
			ContentList: make([]users.Content, 0),
		},
	}
	return u.s.CreateUserProfile(ctx, p)
}

func (u Users) FindUserProfileByID(ctx context.Context, ID string) (*users.UserProfile, error) {
	return u.s.FindUserProfileByID(ctx, ID)
}

func (u Users) RemoveUserProfileByID(ctx context.Context, ID string) error {
	return u.s.RemoveUserProfileByID(ctx, ID)
}
