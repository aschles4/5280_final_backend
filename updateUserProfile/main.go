package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	dbUsers "github.com/aschles4/finalProject/internal/pkg/dynamo/users"
	"github.com/aschles4/finalProject/internal/services/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type UpdateUserProfileRequest struct {
	Password       string                   `json:"password"`
	Name           string                   `json:"name"`
	StreamAccounts []dbUsers.StreamAccounts `json:"streamAccounts"`
}

type UpdateUserProfileResponse struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
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

func (h Handler) HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var token string
	if val, ok := event.Headers["Authorization"]; ok {
		token = strings.Split(val, " ")[1]
	}

	if token == "" {
		return h.handleError(http.StatusBadRequest, nil, "Authorization Header is required")
	}

	js, err := json.Marshal(event.Body)
	if err != nil {
		return h.handleError(http.StatusInternalServerError, err, "Failed to marshal request")
	}
	println(string(js))

	var req UpdateUserProfileRequest
	err = json.Unmarshal([]byte(event.Body), &req)
	if err != nil {
		return h.handleError(http.StatusBadRequest, err, "Failed to marshal event")
	}

	act, err := h.U.FindUserAccountByToken(ctx, token)
	if err != nil {
		return h.handleError(http.StatusUnauthorized, err, "Failed to authorize request")
	}

	prof, err := h.U.FindUserProfileByID(ctx, act.ID)
	if err != nil {
		return h.handleError(http.StatusInternalServerError, err, "Failed to access user profile")
	}

	//Delete User Profile
	err = h.U.RemoveUserProfileByID(ctx, act.ID)
	if err != nil {
		return h.handleError(http.StatusInternalServerError, err, "Failed to Remove Old User Profile")
	}
	//Delete users account
	err = h.U.RemoveUserAccountByID(ctx, act.ID)
	if err != nil {
		return h.handleError(http.StatusInternalServerError, err, "Failed to Remove User Account")
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
		return h.handleError(http.StatusInternalServerError, err, "Failed to Create User Account")
	}

	//create users account
	err = h.U.CreateUserAccountWithToken(ctx, act.ID, act.Email, pass, token)
	if err != nil {
		return h.handleError(http.StatusInternalServerError, err, "Failed to Create User Account")
	}

	//return
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}

func (h Handler) handleError(status int, err error, message string) (events.APIGatewayProxyResponse, error) {
	h.l.Error().Msg(message)
	if err != nil {
		h.l.Error().Msg(err.Error())
	}

	js, err := json.Marshal(UpdateUserProfileResponse{
		Message: message,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "{\"message\":\"InternalServerError\"}",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(js),
	}, nil
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
