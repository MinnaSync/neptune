package animepahe

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func (c *Client) GetEpisodeStreamingLinks(ctx context.Context, animeId int, sessionId string, episodeId string) ([]*EpisodeStreamingLink, error) {
	cachedLinks, _ := c.Cache.GetStreamingLinks(ctx, animeId, episodeId)
	if cachedLinks != nil {
		return cachedLinks, nil
	}

	u := fmt.Sprintf("/play/%s/%s", sessionId, episodeId)
	req, err := c.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	html := new(string)
	_, err = c.Do(ctx, req, html)
	if err != nil {
		return nil, err
	}

	streamingLinks := make([]*EpisodeStreamingLink, 0)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	doc.Find("#resolutionMenu > button").Each(func(i int, s *goquery.Selection) {
		resolution, _ := s.Attr("data-resolution")
		language, _ := s.Attr("data-audio")
		source, _ := s.Attr("data-src")

		file, err := c.kwik.Extract(ctx, source)
		if err != nil {
			log.Error(err.Error())
			return
		}

		streamingLink := EpisodeStreamingLink{
			Resolution: resolution,
			Language:   language,
			URL:        *file,
		}

		streamingLinks = append(streamingLinks, &streamingLink)
	})

	err = c.Cache.SetStreamingLinks(ctx, animeId, episodeId, streamingLinks)
	if err != nil {
		log.Error(err)
	}

	return streamingLinks, nil
}
