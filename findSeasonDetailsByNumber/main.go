package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aschles4/finalProject/internal/services/content"
	"github.com/aschles4/finalProject/internal/services/users"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type findSeasonDetailsByNumberResponse struct {
	Details *content.SeasonDetails `json:"details,omitempty"`
	Status  int                    `json:"status,omitempty"`
	Message string                 `json:"message,omitempty"`
}

type Env struct {
	Connection  string `required:"true" default:"https://dynamodb.us-east-1.amazonaws.com" envconfig:"CONNECTION"`
	Region      string `required:"true" default:"us-east-1" envconfig:"REGION"`
	IMBDKey     string `required:"true" default:"" envconfig:"IMBD_KEY"`
	GuideBoxKey string `required:"true" default:"" envconfig:"GB_KEY"`
}

type Handler struct {
	Env Env
	U   *users.Users
	C   *content.Content
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

	_, err := h.U.FindUserAccountByToken(ctx, token)
	if err != nil {
		return h.handleError(http.StatusUnauthorized, err, "Failed to authorize request")
	}

	var showId string
	if val, ok := event.PathParameters["showid"]; ok {
		showId = val
	}

	if showId == "" {
		return h.handleError(http.StatusBadRequest, nil, "Show ID is required")
	}

	var seasonNumber string
	if val, ok := event.PathParameters["season_num"]; ok {
		seasonNumber = val
	}

	if seasonNumber == "" {
		return h.handleError(http.StatusBadRequest, nil, "Season Number is required")
	}

	d, err := h.C.FindSeasonDetailsByNumber(ctx, showId, seasonNumber)
	if err != nil {
		return h.handleError(http.StatusInternalServerError, err, "Failed to access content details")
	}

	//return
	js, err := json.Marshal(findSeasonDetailsByNumberResponse{
		Details: d,
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

	js, err := json.Marshal(findSeasonDetailsByNumberResponse{
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

	c, err := content.NewContentService(e.IMBDKey)
	if err != nil {
		l.Info().Msg(err.Error())
		l.Fatal().Msg("failed to connect to users service")
	}

	h := Handler{
		l:   l,
		Env: e,
		U:   u,
		C:   c,
	}

	lambda.Start(h.HandleRequest)
}
