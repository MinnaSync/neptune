package animepahe

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) Search(ctx context.Context, term string, page int) (*SearchResponse, *http.Response, error) {
	u := c.buildQuery(url.Values{
		"m": {"search"},
		"q": {term},
		"p": {fmt.Sprintf("%d", page)},
	})

	req, err := c.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	searchResponse := new(SearchResponse)
	resp, err := c.Do(ctx, req, searchResponse)
	if err != nil {
		return nil, resp, err
	}

	return searchResponse, resp, nil
}

func (c *Client) Releases(ctx context.Context, session string, page int) (*EpisodeResponse, *http.Response, error) {
	u := c.buildQuery(url.Values{
		"m":    {"release"},
		"id":   {session},
		"sort": {"episode_desc"},
		"page": {fmt.Sprintf("%d", page)},
	})

	req, err := c.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	episodeResponse := new(EpisodeResponse)
	resp, err := c.Do(ctx, req, episodeResponse)
	if err != nil {
		return nil, resp, err
	}

	return episodeResponse, resp, nil
}
