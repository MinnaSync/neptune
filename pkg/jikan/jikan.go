package jikan

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/redis/go-redis/v9"
)

var (
	BaseURL = "https://api.jikan.moe"
)

type service struct {
	client *Client
}

type Client struct {
	client *http.Client
	common service

	BaseURL *url.URL

	Cache        *JikanCache
	AnimeService *Anime
}

func NewClient(redis *redis.Client) *Client {
	c := &Client{}

	c.client = http.DefaultClient
	c.common.client = c
	c.BaseURL, _ = url.Parse(BaseURL)

	c.AnimeService = (*Anime)(&c.common)

	c.Cache = NewJikanCache(redis)

	return c
}

// Jikan only supports GET requests.
func (c *Client) NewGetRequest(urlStr string) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
