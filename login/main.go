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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token,omitempty"`
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

func (h Handler) HandleRequest(ctx context.Context, req LoginRequest) LoginResponse {
	if req.Email == "" {
		h.l.Info().Msg("Email is required")
		return LoginResponse{
			Status:  http.StatusBadRequest,
			Message: "Email is required",
		}
	}

	if req.Password == "" {
		h.l.Info().Msg("Password is required")
		return LoginResponse{
			Status:  http.StatusBadRequest,
			Message: "Password is required",
		}
	}

	//Login User Here
	token, err := h.U.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		h.l.Error().Msg("Failed to Login User")
		h.l.Error().Msg(err.Error())
		return LoginResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to Login User",
		}
	}

	//return users id & login token
	return LoginResponse{
		Token: *token,
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
