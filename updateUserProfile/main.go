package main

import (
	"context"
	"net/http"
	"os"

	dbUsers "github.com/aschles4/finalProject/internal/pkg/dynamo/users"
	"github.com/aschles4/finalProject/internal/services/users"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type UpdateUserProfileRequest struct {
	Token          string                   `json:"token"`
	Password       string                   `json:"password"`
	Name           string                   `json:"name"`
	StreamAccounts []dbUsers.StreamAccounts `json:"streamAccounts"`
}

type UpdateUserProfileResponse struct {
	Status  int    `json:"token,omitempty"`
	Message string `json:"errorMessage,omitempty"`
}

type Env struct {
	Connection string `required:"true" default:"https://dynamodb.us-east-1.amazonaws.com" envconfig:"CONNECTION"`
	Region     string `required:"true" default:"us-east-1" envconfig:"REGION"`
}

type Handler struct {
	Env Env
	U   *users.Users
	l   zerolog.Logger
}

func (h Handler) HandleRequest(ctx context.Context, req UpdateUserProfileRequest) UpdateUserProfileResponse {
	if req.Token == "" {
		h.l.Info().Msg("Token is required")
		return UpdateUserProfileResponse{
			Status:  http.StatusBadRequest,
			Message: "Token is required",
		}
	}

	act, err := h.U.FindUserAccountByToken(ctx, req.Token)
	if err != nil {
		h.l.Error().Msg("Failed to access user account by token")
		return UpdateUserProfileResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to access user account by token",
		}
	}

	prof, err := h.U.FindUserProfileByID(ctx, act.ID)
	if err != nil {
		h.l.Error().Msg("Failed to access user profile by ID")
		return UpdateUserProfileResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to access user profile",
		}
	}

	//Delete User Profile
	err = h.U.RemoveUserProfileByID(ctx, act.ID)
	if err != nil {
		h.l.Error().Msg("Failed to Remove User Profile")
		h.l.Error().Msg(err.Error())
		return UpdateUserProfileResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to Remove Old User Profile",
		}
	}
	//Delete users account
	err = h.U.RemoveUserAccountByID(ctx, act.ID)
	if err != nil {
		h.l.Error().Msg("Failed to Remove User Account")
		h.l.Error().Msg(err.Error())
		return UpdateUserProfileResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to Create User Account",
		}
	}

	//Update values if they are updated
	pass := req.Password
	if req.Password == "" {
		pass = act.Password
	}

	name := req.Name
	if req.Name == "" {
		pass = prof.Name
	}

	//create users profile
	err = h.U.CreateUserProfile(ctx, act.ID, name, act.Email, req.StreamAccounts)
	if err != nil {
		h.l.Error().Msg("Failed to Create User Profile")
		h.l.Error().Msg(err.Error())
		return UpdateUserProfileResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to Create User Profile",
		}
	}

	//create users account
	err = h.U.CreateUserAccountWithToken(ctx, act.ID, act.Email, pass, req.Token)
	if err != nil {
		h.l.Error().Msg("Failed to Create User Account")
		h.l.Error().Msg(err.Error())
		return UpdateUserProfileResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to Create User Account",
		}
	}

	//return
	return UpdateUserProfileResponse{
		Status: http.StatusNoContent,
	}
}

func main() {
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()

	var e Env
	err := envconfig.Process("", &e)
	if err != nil {
		l.Info().Msg(err.Error())
		l.Fatal().Msg("failed to parse envs")
	}

	u, err := users.NewUsersService(e.Connection, e.Region)
	if err != nil {
		l.Info().Msg(err.Error())
		l.Fatal().Msg("failed to connect to users service")
	}

	h := Handler{
		l:   l,
		Env: e,
		U:   u,
	}

	lambda.Start(h.HandleRequest)
}
