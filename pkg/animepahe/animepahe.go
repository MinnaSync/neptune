package animepahe

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/minna-sync/neptune/config"
	"github.com/minna-sync/neptune/pkg/extractors"
	"github.com/redis/go-redis/v9"
)

var (
	BaseURL = "https://" + config.C.ProviderURLs.Animepahe
)

type HeadersRoundTripper struct {
	http.RoundTripper
}

func (t *HeadersRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())

	headers := map[string]string{
		"Accept":           "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Encoding":  "gzip, deflate, br",
		"Accept-Language":  "en-US,en;q=0.9",
		"Connection":       "keep-alive",
		"Cookie":           "__ddg2_=;",
		"DNT":              "1",
		"Host":             req.URL.Hostname(),
		"Sec-Fetch-Dest":   "empty",
		"Sec-Fetch-Mode":   "cors",
		"X-Requested-With": "XMLHttpRequest",
		"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:143.0) Gecko/20100101 Firefox/143.0",
	}

	for k, v := range headers {
		clone.Header.Set(k, v)
	}

	return t.RoundTripper.RoundTrip(clone)
}

type Client struct {
	client *http.Client
	kwik   *extractors.KwikExtractor

	BaseURL *url.URL
	Cache   *AnimepaheCache
}

func NewClient(redis *redis.Client) *Client {
	c := &Client{}

	c.client = &http.Client{
		Transport: &HeadersRoundTripper{http.DefaultTransport},
	}
	c.kwik = extractors.NewKwikExtractor()
	c.BaseURL, _ = url.Parse(BaseURL)
	c.Cache = NewAnimepaheCache(redis)

	return c
}

func (c *Client) buildQuery(query url.Values) string {
	return fmt.Sprintf("/api?%s", query.Encode())
}

func (c *Client) NewRequest(method, urlStr string, body any) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err

	}
	if body != nil {
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(body)
		if err != nil {
			return nil, err
		}

		req.Body = io.NopCloser(&buf)
		req.Header.Set("Content-Type", "application/json")
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

	if resp.Header.Get("Content-Encoding") == "gzip" {
		contentType := resp.Header.Get("Content-Type")

		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		body, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		switch {
		case strings.HasPrefix(contentType, "text/html"):
			if s, ok := v.(*string); ok {
				*s = string(body)
			} else {
				return resp, fmt.Errorf("expected string, got %T", v)
			}

			return resp, nil
		default:
			err = json.Unmarshal(body, v)
			if err != nil {
				return nil, err
			}
		}
	} else {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}
