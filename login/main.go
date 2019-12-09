package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aschles4/finalProject/internal/services/users"
	"github.com/aws/aws-lambda-go/events"
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
	js, err := json.Marshal(event.Body)
	if err != nil {
		return h.handleError(http.StatusInternalServerError, err, "Failed to marshal request")
	}
	println(string(js))

	var req LoginRequest
	err = json.Unmarshal([]byte(event.Body), &req)
	if err != nil {
		return h.handleError(http.StatusBadRequest, err, "Failed to marshal event")
	}

	if req.Email == "" {
		return h.handleError(http.StatusBadRequest, nil, "Email is required")
	}

	if req.Password == "" {
		return h.handleError(http.StatusBadRequest, nil, "Password is required")
	}

	//Login User Here
	token, err := h.U.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		return h.handleError(http.StatusInternalServerError, err, "Failed to marshal response")
	}

	//return users id & login token
	js, err = json.Marshal(LoginResponse{
		Token: token,
	})
	if err != nil {
		return h.handleError(http.StatusInternalServerError, err, "Failed to marshal response")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func (h Handler) handleError(status int, err error, message string) (events.APIGatewayProxyResponse, error) {
	h.l.Error().Msg(message)
	if err != nil {
		h.l.Error().Msg(err.Error())
	}

	js, err := json.Marshal(LoginResponse{
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
