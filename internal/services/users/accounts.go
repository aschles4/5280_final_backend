package users

import (
	"context"

	"github.com/aschles4/finalProject/internal/pkg/dynamo/users"
	"github.com/segmentio/ksuid"
)

func (u Users) CreateUserAccount(ctx context.Context, ID, email, pass string) error {

	p, err := u.passEncrypt(ctx, pass)
	if err != nil {
		return err
	}

	a := users.UserAccount{
		ID:       ID,
		Email:    email,
		Password: p,
	}
	return u.s.CreateUserAccount(ctx, a)
}

func (u Users) CreateUserAccountWithToken(ctx context.Context, ID, email, pass, token string) error {

	p, err := u.passEncrypt(ctx, pass)
	if err != nil {
		return err
	}

	a := users.UserAccount{
		ID:       ID,
		Email:    email,
		Password: p,
		Token:    token,
	}
	return u.s.CreateUserAccount(ctx, a)
}

func (u Users) FindUserAccountByEmail(ctx context.Context, email string) (*users.UserAccount, error) {
	return u.s.FindUserAccountByEmail(ctx, email)
}

func (u Users) FindUserAccountByToken(ctx context.Context, token string) (*users.UserAccount, error) {
	return u.s.FindUserAccountByToken(ctx, token)
}

func (u Users) LoginUser(ctx context.Context, email, pass string) (*string, error) {
	//find account
	act, err := u.FindUserAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	token := u.generateToken(ctx)
	//update account
	_, err = u.s.UpdateUserAccountToken(ctx, act.ID, token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (u Users) LogoutUserByID(ctx context.Context, ID string) error {
	_, err := u.s.UpdateUserAccountToken(ctx, ID, "")
	if err != nil {
		return err
	}

	return nil
}

func (u Users) RemoveUserAccountByID(ctx context.Context, ID string) error {
	return u.s.RemoveUserAccountByID(ctx, ID)
}

//Helpers

func (u Users) passEncrypt(ctx context.Context, pass string) (string, error) {
	//TODO ENCRYPT PASSWORD!!
	return pass, nil
}

func (u Users) generateToken(ctx context.Context) string {
	return ksuid.New().String()
}
