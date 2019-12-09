package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/aschles4/finalProject/internal/services/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type LogoutResponse struct {
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

	act, err := h.U.FindUserAccountByToken(ctx, token)
	if err != nil {
		return h.handleError(http.StatusUnauthorized, err, "Failed to authorize request")
	}

	//Logout User Here
	err = h.U.LogoutUserByID(ctx, act.ID)
	if err != nil {
		return h.handleError(http.StatusInternalServerError, err, "Failed to Logout User")
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

	js, err := json.Marshal(LogoutResponse{
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
