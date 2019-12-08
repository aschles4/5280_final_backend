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
	"github.com/segmentio/ksuid"
)

type SignupRequest struct {
	Email          string                   `json:"email"`
	Password       string                   `json:"password"`
	Name           string                   `json:"name"`
	StreamAccounts []dbUsers.StreamAccounts `json:"streamAccounts"`
}

type SignupResponse struct {
	UserID  string `json:"userId,omitempty"`
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

func (h Handler) HandleRequest(ctx context.Context, req SignupRequest) SignupResponse {
	if req.Email == "" {
		h.l.Info().Msg("Email is required")
		return SignupResponse{
			Status:  http.StatusBadRequest,
			Message: "Email is required",
		}
	}

	if req.Password == "" {
		h.l.Info().Msg("Password is required")
		return SignupResponse{
			Status:  http.StatusBadRequest,
			Message: "Password is required",
		}
	}

	if req.Name == "" {
		h.l.Info().Msg("Name is required")
		return SignupResponse{
			Status:  http.StatusBadRequest,
			Message: "Name is required",
		}
	}

	act, err := h.U.FindUserAccountByEmail(ctx, req.Email)
	if act != nil {
		h.l.Error().Msg("Account already created with email")
		return SignupResponse{
			Status:  http.StatusBadRequest,
			Message: "Account already created with email",
		}
	}

	//create users profile
	ID := ksuid.New().String()
	err = h.U.CreateUserProfile(ctx, ID, req.Name, req.Email, req.StreamAccounts)
	if err != nil {
		h.l.Error().Msg("Failed to Create User Profile")
		h.l.Error().Msg(err.Error())
		return SignupResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to Create User Profile",
		}
	}

	//create users account
	err = h.U.CreateUserAccount(ctx, ID, req.Email, req.Password)
	if err != nil {
		h.l.Error().Msg("Failed to Create User Account")
		h.l.Error().Msg(err.Error())
		return SignupResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to Create User Account",
		}
	}

	//Login User Here
	token, err := h.U.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		h.l.Error().Msg("Failed to Login User On SignUp")
		h.l.Error().Msg(err.Error())
		return SignupResponse{
			UserID: ID,
		}
	}

	//return users id & login token
	return SignupResponse{
		UserID: ID,
		Token:  *token,
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
