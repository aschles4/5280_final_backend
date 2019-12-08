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

type findEpisodeDetailsByNumberRequest struct {
	Token  string `json:"token"`
	Number int    `json:"number"`
}

type findEpisodeDetailsByNumberResponse struct {
	Details *content.EpisodeDetails `json:"details,omitempty"`
	Status  int                     `json:"token,omitempty"`
	Message string                  `json:"errorMessage,omitempty"`
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

func (h Handler) HandleRequest(ctx context.Context, req findEpisodeDetailsByNumberRequest) findEpisodeDetailsByNumberResponse {
	if req.Token == "" {
		h.l.Info().Msg("Token is required")
		return findEpisodeDetailsByNumberResponse{
			Status:  http.StatusBadRequest,
			Message: "Token is required",
		}
	}

	_, err := h.U.FindUserAccountByToken(ctx, req.Token)
	if err != nil {
		h.l.Error().Msg("Failed to access user account by token")
		return findEpisodeDetailsByNumberResponse{
			Status:  http.StatusUnauthorized,
			Message: "Failed to access user account by token",
		}
	}

	d, err := h.C.FindEpisodeDetailsByNumber(ctx, req.Number)
	if err != nil {
		h.l.Error().Msg("Failed to access content details")
		h.l.Error().Msg(err.Error())
		return findEpisodeDetailsByNumberResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to access content details",
		}
	}

	//return
	return findEpisodeDetailsByNumberResponse{
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
