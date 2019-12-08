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

type FindUserProfileRequest struct {
	Token string `json:"token"`
}

type FindUserProfileResponse struct {
	Email          string                   `json:"email,omitempty"`
	Name           string                   `json:"name,omitempty"`
	StreamAccounts []dbUsers.StreamAccounts `json:"streamAccounts,omitempty"`
	Status         int                      `json:"token,omitempty"`
	Message        string                   `json:"errorMessage,omitempty"`
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

func (h Handler) HandleRequest(ctx context.Context, req FindUserProfileRequest) FindUserProfileResponse {
	if req.Token == "" {
		h.l.Info().Msg("Token is required")
		return FindUserProfileResponse{
			Status:  http.StatusBadRequest,
			Message: "Token is required",
		}
	}

	act, err := h.U.FindUserAccountByToken(ctx, req.Token)
	if err != nil {
		h.l.Error().Msg("Failed to access user account by token")
		return FindUserProfileResponse{
			Status:  http.StatusUnauthorized,
			Message: "Failed to access user account by token",
		}
	}

	prof, err := h.U.FindUserProfileByID(ctx, act.ID)
	if err != nil {
		h.l.Error().Msg("Failed to access user profile by ID")
		return FindUserProfileResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to access user profile",
		}
	}

	//return
	return FindUserProfileResponse{
		Email:          act.Email,
		Name:           prof.Name,
		StreamAccounts: prof.StreamAccounts,
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
