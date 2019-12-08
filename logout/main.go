package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aschles4/finalProject/internal/services/users"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type LogoutRequest struct {
	Token string `json:"token"`
}

type LogoutResponse struct {
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

func (h Handler) HandleRequest(ctx context.Context, req LogoutRequest) LogoutResponse {
	if req.Token == "" {
		h.l.Info().Msg("Token is required")
		return LogoutResponse{
			Status:  http.StatusBadRequest,
			Message: "Token is required",
		}
	}

	act, err := h.U.FindUserAccountByToken(ctx, req.Token)
	if err != nil {
		h.l.Error().Msg("Failed to access user account by token")
		return LogoutResponse{
			Status:  http.StatusUnauthorized,
			Message: "Failed to access user account by token",
		}
	}

	//Logout User Here
	err := h.U.LogoutUserByID(ctx, act.ID)
	if err != nil {
		h.l.Error().Msg("Failed to Logout User")
		h.l.Error().Msg(err.Error())
		return LogoutResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to Logout User",
		}
	}

	//return
	return LogoutResponse{
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
