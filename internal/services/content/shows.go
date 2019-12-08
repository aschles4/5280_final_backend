package content

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ShowDetails struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	URL         string   `json:"thumbnailURL"`
	Seasons     []Season `json:"seasons"`
}

type SeasonDetails struct {
	SeasonNumber int       `json:"season_number"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	URL          string    `json:"thumbnailURL"`
	Episodes     []Episode `json:"episodes"`
}

type EpisodeDetails struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	URL         string     `json:"thumbnailURL"`
	WatchNow    []WatchNow `json:"watchNow"`
}

type Season struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PosterPath   string `json:"thumbnailURL"`
	SeasonNumber int    `json:"season_number"`
}

type Episode struct {
	EpisodeNumber int    `json:"episode_number"`
	PosterPath    string `json:"thumbnailURL"`
}

type ResponseSeason struct {
	AirDate      string `json:"air_date"`
	EpisodeCount int    `json:"episode_count"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	SeasonNumber int    `json:"season_number"`
}

type ShowDetailsResponse struct {
	BackdropPath   string        `json:"backdrop_path"`
	CreatedBy      []interface{} `json:"created_by"`
	EpisodeRunTime []int         `json:"episode_run_time"`
	FirstAirDate   string        `json:"first_air_date"`
	Genres         []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage         string   `json:"homepage"`
	ID               int      `json:"id"`
	InProduction     bool     `json:"in_production"`
	Languages        []string `json:"languages"`
	LastAirDate      string   `json:"last_air_date"`
	LastEpisodeToAir struct {
		AirDate        string      `json:"air_date"`
		EpisodeNumber  int         `json:"episode_number"`
		ID             int         `json:"id"`
		Name           string      `json:"name"`
		Overview       string      `json:"overview"`
		ProductionCode string      `json:"production_code"`
		SeasonNumber   int         `json:"season_number"`
		ShowID         int         `json:"show_id"`
		StillPath      interface{} `json:"still_path"`
		VoteAverage    int         `json:"vote_average"`
		VoteCount      int         `json:"vote_count"`
	} `json:"last_episode_to_air"`
	Name             string      `json:"name"`
	NextEpisodeToAir interface{} `json:"next_episode_to_air"`
	Networks         []struct {
		Name          string `json:"name"`
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		OriginCountry string `json:"origin_country"`
	} `json:"networks"`
	NumberOfEpisodes    int      `json:"number_of_episodes"`
	NumberOfSeasons     int      `json:"number_of_seasons"`
	OriginCountry       []string `json:"origin_country"`
	OriginalLanguage    string   `json:"original_language"`
	OriginalName        string   `json:"original_name"`
	Overview            string   `json:"overview"`
	Popularity          float64  `json:"popularity"`
	PosterPath          string   `json:"poster_path"`
	ProductionCompanies []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	Seasons     []ResponseSeason `json:"seasons"`
	Status      string           `json:"status"`
	Type        string           `json:"type"`
	VoteAverage float64          `json:"vote_average"`
	VoteCount   int              `json:"vote_count"`
}

type SeasonDetailsResponse struct {
	RequestID string `json:"_id"`
	AirDate   string `json:"air_date"`
	Episodes  []struct {
		AirDate        string        `json:"air_date"`
		EpisodeNumber  int           `json:"episode_number"`
		ID             int           `json:"id"`
		Name           string        `json:"name"`
		Overview       string        `json:"overview"`
		ProductionCode string        `json:"production_code"`
		SeasonNumber   int           `json:"season_number"`
		ShowID         int           `json:"show_id"`
		StillPath      string        `json:"still_path"`
		VoteAverage    int           `json:"vote_average"`
		VoteCount      int           `json:"vote_count"`
		Crew           []interface{} `json:"crew"`
		GuestStars     []interface{} `json:"guest_stars"`
	} `json:"episodes"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	ID           int    `json:"id"`
	PosterPath   string `json:"poster_path"`
	SeasonNumber int    `json:"season_number"`
}

type EpisodeDetailResponse struct {
	AirDate string `json:"air_date"`
	Crew    []struct {
		ID          int    `json:"id"`
		CreditID    string `json:"credit_id"`
		Name        string `json:"name"`
		Department  string `json:"department"`
		Job         string `json:"job"`
		ProfilePath string `json:"profile_path"`
	} `json:"crew"`
	EpisodeNumber int `json:"episode_number"`
	GuestStars    []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		CreditID    string `json:"credit_id"`
		Character   string `json:"character"`
		Order       int    `json:"order"`
		ProfilePath string `json:"profile_path"`
	} `json:"guest_stars"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ID             int     `json:"id"`
	ProductionCode string  `json:"production_code"`
	SeasonNumber   int     `json:"season_number"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
}

func (c Content) FindShowSuggestions(ctx context.Context) (*[]Suggestion, error) {
	query := []string{"Action", "Comedy", "Family", "Science Fiction"}
	suggestions := make([]Suggestion, 0)
	for _, q := range query {
		url := fmt.Sprintf("https://api.themoviedb.org/3/search/tv?include_adult=false&page=1&query=%s&language=en-US&api_key=%s", q, c.ApiKey)
		payload := strings.NewReader("{}")

		req, _ := http.NewRequest("GET", url, payload)

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var resp MovieSearchResponse
		err := json.Unmarshal(body, &resp)
		if err != nil {
			return nil, err
		}

		thumbnails := make([]Thumbnail, 0)
		for _, r := range resp.Results {
			thumbnails = append(thumbnails, Thumbnail{
				ID:  fmt.Sprintf("%v", r.ID),
				URL: fmt.Sprintf("https://image.tmdb.org/t/p/w185%v", r.PosterPath),
			})
		}

		suggestions = append(suggestions, Suggestion{
			Type:     "TV",
			Category: q,
			List:     thumbnails,
		})
	}

	return &suggestions, nil
}

func (c Content) FindShowDetailsByID(ctx context.Context, ID string) (*ShowDetails, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/tv/%s?language=en-US&api_key=%s", ID, c.ApiKey)

	payload := strings.NewReader("{}")
	req, _ := http.NewRequest("GET", url, payload)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var resp ShowDetailsResponse

	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	seasons := make([]Season, 0)
	for _, s := range resp.Seasons {
		seasons = append(seasons, Season{
			ID:           fmt.Sprintf("%v", s.ID),
			Title:        s.Name,
			Description:  s.Overview,
			PosterPath:   fmt.Sprintf("https://image.tmdb.org/t/p/w185%v", s.PosterPath),
			SeasonNumber: s.SeasonNumber,
		})
	}

	d := ShowDetails{
		ID:          fmt.Sprintf("%v", resp.ID),
		Title:       fmt.Sprintf("%v", resp.Name),
		Description: fmt.Sprintf("%v", resp.Overview),
		URL:         fmt.Sprintf("https://image.tmdb.org/t/p/w185%v", resp.PosterPath),
		Seasons:     seasons,
	}

	return &d, nil
}

func (c Content) FindSeasonDetailsByNumber(ctx context.Context, number int) (*SeasonDetails, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/tv/season/%d?language=en-US&api_key=%s", number, c.ApiKey)

	payload := strings.NewReader("{}")
	req, _ := http.NewRequest("GET", url, payload)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var resp SeasonDetailsResponse

	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	episodes := make([]Episode, 0)
	for _, e := range resp.Episodes {
		episodes = append(episodes, Episode{
			PosterPath:    fmt.Sprintf("https://image.tmdb.org/t/p/w185%v", e.StillPath),
			EpisodeNumber: e.EpisodeNumber,
		})
	}

	d := SeasonDetails{
		SeasonNumber: resp.SeasonNumber,
		Title:        fmt.Sprintf("%v", resp.Name),
		Description:  fmt.Sprintf("%v", resp.Overview),
		URL:          fmt.Sprintf("https://image.tmdb.org/t/p/w185%v", resp.PosterPath),
		Episodes:     episodes,
	}

	return &d, nil
}

func (c Content) FindEpisodeDetailsByNumber(ctx context.Context, number int) (*EpisodeDetails, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/tv/season/episode/%d?language=en-US&api_key=%s", number, c.ApiKey)

	payload := strings.NewReader("{}")
	req, _ := http.NewRequest("GET", url, payload)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var resp EpisodeDetailResponse

	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	d := EpisodeDetails{
		ID:          fmt.Sprintf("%v", resp.ID),
		Title:       fmt.Sprintf("%v", resp.Name),
		Description: fmt.Sprintf("%v", resp.Overview),
		URL:         fmt.Sprintf("https://image.tmdb.org/t/p/w185%v", resp.StillPath),
	}

	return &d, nil
}
