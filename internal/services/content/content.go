package content

import "context"

type Content struct {
	ApiKey string `json:"key"`
}

type Thumbnail struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Suggestion struct {
	Type     string      `json:"type"`
	Category string      `json:"category"`
	List     []Thumbnail `json:"list"`
}

func NewContentService(key string) (*Content, error) {
	return &Content{
		ApiKey: key,
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
