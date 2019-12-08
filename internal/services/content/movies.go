package content

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type WatchNow struct {
	ID  string `json:"serviceId"`
	URL string `json:"url"`
}

type MovieDetails struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	URL         string     `json:"thumbnailURL"`
	WatchNow    []WatchNow `json:"watchNow"`
}

type MovieSearchResponse struct {
	Page         int       `json:"page"`
	TotalResults int       `json:"total_results"`
	TotalPages   int       `json:"total_pages"`
	Results      []Results `json:"results"`
}

type Results struct {
	Popularity       float64 `json:"popularity"`
	VoteCount        int     `json:"vote_count"`
	Video            bool    `json:"video"`
	PosterPath       string  `json:"poster_path"`
	ID               int     `json:"id"`
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	GenreIds         []int   `json:"genre_ids"`
	Title            string  `json:"title"`
	VoteAverage      float64 `json:"vote_average"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
}

func (c Content) FindMovieSuggestions(ctx context.Context) (*[]Suggestion, error) {
	query := []string{"Action", "Comedy", "Family", "Science Fiction"}
	suggestions := make([]Suggestion, 0)
	for _, q := range query {
		url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?include_adult=false&page=1&query=%s&language=en-US&api_key=%s", q, c.ApiKey)
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
			Type:     "Movie",
			Category: q,
			List:     thumbnails,
		})
	}

	return &suggestions, nil
}

func (c Content) FindMovieDetailsByID(ctx context.Context, ID string) (*MovieDetails, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?language=en-US&api_key=%s", ID, c.ApiKey)

	payload := strings.NewReader("{}")
	req, _ := http.NewRequest("GET", url, payload)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var b map[string]interface{}

	err := json.Unmarshal(body, &b)
	if err != nil {
		return nil, err
	}

	var description string
	if val, ok := b["overview"]; ok {
		description = fmt.Sprintf("%v", val)
	}

	var title string
	if val, ok := b["title"]; ok {
		title = fmt.Sprintf("%v", val)
	}

	var u string
	if val, ok := b["poster_path"]; ok {
		url = fmt.Sprintf("https://image.tmdb.org/t/p/w185%v", val)
	}

	d := MovieDetails{
		ID:          ID,
		Title:       title,
		Description: description,
		URL:         u,
	}

	return &d, nil
}
