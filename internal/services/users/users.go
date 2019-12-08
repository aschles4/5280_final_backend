package users

import (
	"github.com/aschles4/finalProject/internal/pkg/dynamo/users"
)

type Users struct {
	s store
}

func NewUsersService(conn, region string) (*Users, error) {
	var s store
	s, err := users.NewStore(conn, region)
	if err != nil {
		return nil, err
	}

	return &Users{
		s: s,
	}, nil
}
