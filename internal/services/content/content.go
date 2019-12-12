package content

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aschles4/finalProject/internal/services/guidebox"
)

type Content struct {
	ApiKey   string            `json:"key"`
	GuideBox guidebox.GuideBox `json:"GuideBox"`
}

type Thumbnail struct {
	ID   string `json:"id"`
	Type string `json:"type,omitempty"`
	URL  string `json:"url"`
}

type Suggestion struct {
	Type     string      `json:"type,omitempty"`
	Category string      `json:"category,omitempty"`
	List     []Thumbnail `json:"list"`
}

type MultiSearchResponse struct {
	Page    int `json:"page"`
	Results []struct {
		PosterPath       interface{}   `json:"poster_path,omitempty"`
		Popularity       float64       `json:"popularity"`
		ID               int           `json:"id"`
		Overview         string        `json:"overview,omitempty"`
		BackdropPath     interface{}   `json:"backdrop_path,omitempty"`
		VoteAverage      float64       `json:"vote_average,omitempty"`
		MediaType        string        `json:"media_type"`
		FirstAirDate     string        `json:"first_air_date,omitempty"`
		OriginCountry    []string      `json:"origin_country,omitempty"`
		GenreIds         []interface{} `json:"genre_ids,omitempty"`
		OriginalLanguage string        `json:"original_language,omitempty"`
		VoteCount        float64       `json:"vote_count,omitempty"`
		Name             string        `json:"name,omitempty"`
		OriginalName     string        `json:"original_name,omitempty"`
		Adult            bool          `json:"adult,omitempty"`
		ReleaseDate      string        `json:"release_date,omitempty"`
		OriginalTitle    string        `json:"original_title,omitempty"`
		Title            string        `json:"title,omitempty"`
		Video            bool          `json:"video,omitempty"`
		ProfilePath      string        `json:"profile_path,omitempty"`
		KnownFor         []struct {
			PosterPath       string  `json:"poster_path"`
			Adult            bool    `json:"adult"`
			Overview         string  `json:"overview"`
			ReleaseDate      string  `json:"release_date"`
			OriginalTitle    string  `json:"original_title"`
			GenreIds         []int   `json:"genre_ids"`
			ID               int     `json:"id"`
			MediaType        string  `json:"media_type"`
			OriginalLanguage string  `json:"original_language"`
			Title            string  `json:"title"`
			BackdropPath     string  `json:"backdrop_path"`
			Popularity       float64 `json:"popularity"`
			VoteCount        float64 `json:"vote_count"`
			Video            bool    `json:"video"`
			VoteAverage      float64 `json:"vote_average"`
		} `json:"known_for,omitempty"`
	} `json:"results"`
	TotalResults int `json:"total_results"`
	TotalPages   int `json:"total_pages"`
}

func NewContentService(tmdbKey, guideboxKey string) (*Content, error) {

	g, err := guidebox.NewGuideBoxService(guideboxKey)
	if err != nil {
		return nil, nil
	}
	return &Content{
		ApiKey:   tmdbKey,
		GuideBox: *g,
	}, nil
}

func (c Content) FindAllContentSuggestions(ctx context.Context) (*[]Suggestion, error) {
	movies, err := c.FindMovieSuggestions(ctx)
	if err != nil {
		return nil, err
	}

	shows, err := c.FindShowSuggestions(ctx)
	if err != nil {
		return nil, err
	}

	s := make([]Suggestion, 0)
	s = append(s, *movies...)
	s = append(s, *shows...)

	return &s, nil
}

func (c Content) Search(ctx context.Context, query string) (*[]Suggestion, error) {
	suggestions := make([]Suggestion, 0)
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/multi?include_adult=false&page=1&query=%s&language=en-US&api_key=%s", query, c.ApiKey)
	payload := strings.NewReader("{}")

	req, _ := http.NewRequest("GET", url, payload)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resp MultiSearchResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	j, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	println(string(j))

	thumbnails := make([]Thumbnail, 0)
	for _, r := range resp.Results {

		var u string
		if r.PosterPath != nil {
			u = fmt.Sprintf("https://image.tmdb.org/t/p/w185%v", r.PosterPath)
		}

		thumbnails = append(thumbnails, Thumbnail{
			ID:   fmt.Sprintf("%d", r.ID),
			Type: r.MediaType,
			URL:  u,
		})
	}

	suggestions = append(suggestions, Suggestion{
		List: thumbnails,
	})

	return &suggestions, nil
}
