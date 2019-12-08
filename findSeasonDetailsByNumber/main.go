package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aschles4/finalProject/internal/services/content"
	"github.com/aschles4/finalProject/internal/services/users"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type findSeasonDetailsByNumberRequest struct {
	Token  string `json:"token"`
	Number int    `json:"number"`
}

type findSeasonDetailsByNumberResponse struct {
	Details *content.SeasonDetails `json:"details,omitempty"`
	Status  int                    `json:"token,omitempty"`
	Message string                 `json:"errorMessage,omitempty"`
}

type Env struct {
	Connection string `required:"true" default:"https://dynamodb.us-east-1.amazonaws.com" envconfig:"CONNECTION"`
	Region     string `required:"true" default:"us-east-1" envconfig:"REGION"`
	IMBDKey    string `required:"true" default:"" envconfig:"IMBD_KEY"`
}

type Handler struct {
	Env Env
	U   *users.Users
	C   *content.Content
	l   zerolog.Logger
}

func (h Handler) HandleRequest(ctx context.Context, req findSeasonDetailsByNumberRequest) findSeasonDetailsByNumberResponse {
	if req.Token == "" {
		h.l.Info().Msg("Token is required")
		return findSeasonDetailsByNumberResponse{
			Status:  http.StatusBadRequest,
			Message: "Token is required",
		}
	}

	_, err := h.U.FindUserAccountByToken(ctx, req.Token)
	if err != nil {
		h.l.Error().Msg("Failed to access user account by token")
		return findSeasonDetailsByNumberResponse{
			Status:  http.StatusUnauthorized,
			Message: "Failed to access user account by token",
		}
	}

	d, err := h.C.FindSeasonDetailsByNumber(ctx, req.Number)
	if err != nil {
		h.l.Error().Msg("Failed to access content details")
		h.l.Error().Msg(err.Error())
		return findSeasonDetailsByNumberResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to access content details",
		}
	}

	//return
	return findSeasonDetailsByNumberResponse{
		Details: d,
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
