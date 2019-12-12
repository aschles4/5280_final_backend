package guidebox

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SeasonSources struct {
	Season         string          `json:"season,omitempty"`
	EpisodeSources []EpisodeSource `json:"episodeSources,omitempty"`
}

type EpisodeSource struct {
	ID             int      `json:"id"`
	TMDBID         int      `json:"tmdb_id"`
	EpisodeNumber  int      `json:"episodeNumber"`
	EpisodeSources []Source `json:"sources,omitempty"`
}

type Source struct {
	Source          string `json:"source,omitempty"`
	DisplayName     string `json:"display_name,omitempty"`
	TvChannel       string `json:"tv_channel,omitempty"`
	Link            string `json:"link,omitempty"`
	AppName         string `json:"app_name,omitempty"`
	AppLink         int    `json:"app_link,omitempty"`
	AppRequired     int    `json:"app_required,omitempty"`
	AppDownloadLink string `json:"app_download_link,omitempty"`
	Formats         []struct {
		Price    string `json:"price,omitempty"`
		Format   string `json:"format,omitempty"`
		Type     string `json:"type,omitempty"`
		PreOrder bool   `json:"pre_order,omitempty"`
	} `json:"formats,omitempty"`
}
type MovieDetailsResponse struct {
	ID               int      `json:"id"`
	Title            string   `json:"title"`
	ReleaseYear      int      `json:"release_year"`
	Themoviedb       int      `json:"themoviedb"`
	OriginalTitle    string   `json:"original_title"`
	AlternateTitles  []string `json:"alternate_titles"`
	Imdb             string   `json:"imdb"`
	PreOrder         bool     `json:"pre_order"`
	InTheaters       bool     `json:"in_theaters"`
	ReleaseDate      string   `json:"release_date"`
	Rating           string   `json:"rating"`
	Rottentomatoes   int      `json:"rottentomatoes"`
	Freebase         string   `json:"freebase"`
	WikipediaID      int      `json:"wikipedia_id"`
	Metacritic       string   `json:"metacritic"`
	CommonSenseMedia string   `json:"common_sense_media"`
	Overview         string   `json:"overview"`
	Poster120X171    string   `json:"poster_120x171"`
	Poster240X342    string   `json:"poster_240x342"`
	Poster400X570    string   `json:"poster_400x570"`
	Social           struct {
		Facebook struct {
			FacebookID int64  `json:"facebook_id"`
			Link       string `json:"link"`
		} `json:"facebook"`
	} `json:"social"`
	Genres []struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"genres"`
	Tags []struct {
		ID  int    `json:"id"`
		Tag string `json:"tag"`
	} `json:"tags"`
	Duration int `json:"duration"`
	Trailers struct {
		Web []struct {
			Type        string `json:"type"`
			Source      string `json:"source"`
			DisplayName string `json:"display_name"`
			Link        string `json:"link"`
			Embed       string `json:"embed"`
		} `json:"web"`
		Ios []struct {
			Type        string `json:"type"`
			Source      string `json:"source"`
			DisplayName string `json:"display_name"`
			Link        string `json:"link"`
			Embed       string `json:"embed"`
		} `json:"ios"`
		Android []struct {
			Type        string `json:"type"`
			Source      string `json:"source"`
			DisplayName string `json:"display_name"`
			Link        string `json:"link"`
			Embed       string `json:"embed"`
		} `json:"android"`
	} `json:"trailers"`
	Writers []struct {
		ID    int         `json:"id"`
		Name  string      `json:"name"`
		Imdb  string      `json:"imdb"`
		Image interface{} `json:"image"`
	} `json:"writers"`
	Directors []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Imdb  string `json:"imdb"`
		Image string `json:"image"`
	} `json:"directors"`
	Cast []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		CharacterName string `json:"character_name"`
		Imdb          string `json:"imdb"`
		Image         string `json:"image"`
	} `json:"cast"`
	FreeWebSources             []Source `json:"free_web_sources"`
	FreeIosSources             []Source `json:"free_ios_sources"`
	FreeAndroidSources         []Source `json:"free_android_sources"`
	TvEverywhereWebSources     []Source `json:"tv_everywhere_web_sources"`
	TvEverywhereIosSources     []Source `json:"tv_everywhere_ios_sources"`
	TvEverywhereAndroidSources []Source `json:"tv_everywhere_android_sources"`
	SubscriptionWebSources     []Source `json:"subscription_web_sources"`
	SubscriptionIosSources     []Source `json:"subscription_ios_sources"`
	SubscriptionAndroidSources []Source `json:"subscription_android_sources"`
	PurchaseWebSources         []Source `json:"purchase_web_sources"`
	PurchaseIosSources         []Source `json:"purchase_ios_sources"`
	PurchaseAndroidSources     []Source `json:"purchase_android_sources"`
	OtherSources               []Source `json:"other_sources"`
}

type EpisodeDetialsResponse struct {
	TotalResults  int `json:"total_results"`
	TotalReturned int `json:"total_returned"`
	Results       []struct {
		ID                         int           `json:"id"`
		Tvdb                       int           `json:"tvdb"`
		ContentType                string        `json:"content_type"`
		IsShadow                   int           `json:"is_shadow"`
		AlternateTvdb              []interface{} `json:"alternate_tvdb"`
		ImdbID                     string        `json:"imdb_id"`
		Themoviedb                 int           `json:"themoviedb"`
		ShowID                     int           `json:"show_id"`
		SeasonNumber               int           `json:"season_number"`
		EpisodeNumber              int           `json:"episode_number"`
		Special                    int           `json:"special"`
		FirstAired                 string        `json:"first_aired"`
		Title                      string        `json:"title"`
		OriginalTitle              string        `json:"original_title"`
		AlternateTitles            []interface{} `json:"alternate_titles"`
		Overview                   string        `json:"overview"`
		Duration                   int           `json:"duration"`
		ProductionCode             string        `json:"production_code"`
		Thumbnail208X117           string        `json:"thumbnail_208x117"`
		Thumbnail304X171           string        `json:"thumbnail_304x171"`
		Thumbnail400X225           string        `json:"thumbnail_400x225"`
		Thumbnail608X342           string        `json:"thumbnail_608x342"`
		FreeWebSources             []Source      `json:"free_web_sources"`
		FreeIosSources             []Source      `json:"free_ios_sources"`
		FreeAndroidSources         []Source      `json:"free_android_sources"`
		TvEverywhereWebSources     []Source      `json:"tv_everywhere_web_sources"`
		TvEverywhereIosSources     []Source      `json:"tv_everywhere_ios_sources"`
		TvEverywhereAndroidSources []Source      `json:"tv_everywhere_android_sources"`
		SubscriptionWebSources     []Source      `json:"subscription_web_sources"`
		SubscriptionIosSources     []Source      `json:"subscription_ios_sources"`
		SubscriptionAndroidSources []Source      `json:"subscription_android_sources"`
		PurchaseWebSources         []Source      `json:"purchase_web_sources"`
		PurchaseIosSources         []Source      `json:"purchase_ios_sources"`
		PurchaseAndroidSources     []Source      `json:"purchase_android_sources"`
	} `json:"results"`
}

type SearchByIDResponse struct {
	ID              int           `json:"id"`
	Title           string        `json:"title"`
	AlternateTitles []interface{} `json:"alternate_titles"`
	ContainerShow   int           `json:"container_show"`
	FirstAired      string        `json:"first_aired"`
	ImdbID          string        `json:"imdb_id"`
	Tvdb            int           `json:"tvdb"`
	Themoviedb      int           `json:"themoviedb"`
	Freebase        string        `json:"freebase"`
	WikipediaID     int           `json:"wikipedia_id"`
	Tvrage          struct {
		TvrageID int    `json:"tvrage_id"`
		Link     string `json:"link"`
	} `json:"tvrage"`
	Artwork208X117 string `json:"artwork_208x117"`
	Artwork304X171 string `json:"artwork_304x171"`
	Artwork448X252 string `json:"artwork_448x252"`
	Artwork608X342 string `json:"artwork_608x342"`
}

type GuideBox struct {
	ApiKey string `json:"key"`
}

func NewGuideBoxService(key string) (*GuideBox, error) {
	return &GuideBox{
		ApiKey: key,
	}, nil
}

func (g GuideBox) SearchByIDAndType(ctx context.Context, t, q string) (string, error) {
	url := fmt.Sprintf("http://api-public.guidebox.com/v2/search?api_key=%s&type=%s&field=id&id_type=themoviedb&query=%s", g.ApiKey, t, q)
	payload := strings.NewReader("{}")

	req, _ := http.NewRequest("GET", url, payload)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resp SearchByIDResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", resp.ID), nil
}

func (g GuideBox) FindMovieSources(ctx context.Context, ID string) (*[]Source, error) {
	url := fmt.Sprintf("http://api-public.guidebox.com/v2/movies/%s?api_key=%s", ID, g.ApiKey)
	payload := strings.NewReader("{}")

	req, _ := http.NewRequest("GET", url, payload)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resp MovieDetailsResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	sources := make([]Source, 0)
	sources = append(sources, resp.FreeAndroidSources...)
	sources = append(sources, resp.SubscriptionAndroidSources...)
	sources = append(sources, resp.PurchaseAndroidSources...)

	return &sources, nil
}

func (g GuideBox) FindEpisodeDetails(ctx context.Context, ID, seasonNum string) (*SeasonSources, error) {
	url := fmt.Sprintf("http://api-public.guidebox.com/v2/shows/%s/episodes?api_key=%s&include_links=android&season=%s", ID, g.ApiKey, seasonNum)
	payload := strings.NewReader("{}")

	req, _ := http.NewRequest("GET", url, payload)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resp EpisodeDetialsResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	epSources := make([]EpisodeSource, 0)

	for _, res := range resp.Results {

		sources := make([]Source, 0)
		sources = append(sources, res.FreeAndroidSources...)
		sources = append(sources, res.SubscriptionAndroidSources...)
		sources = append(sources, res.PurchaseAndroidSources...)

		epSources = append(epSources, EpisodeSource{
			ID:             res.ID,
			TMDBID:         res.Themoviedb,
			EpisodeNumber:  res.EpisodeNumber,
			EpisodeSources: sources,
		})
	}

	return &SeasonSources{
		Season:         seasonNum,
		EpisodeSources: epSources,
	}, nil
}
