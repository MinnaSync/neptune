package animepahe

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/minna-sync/neptune/pkg/jikan"
)

func (c *Client) stringCompare(a, b string) float64 {
	distance := levenshtein.ComputeDistance(a, b)
	return (1.0 - float64(distance)/float64(max(len(b), len(a)))) * 100
}

func (c *Client) GetMatchedResult(ctx context.Context, info *jikan.AnimeInfoBase) (*string, error) {
	cachedSession, _ := c.Cache.GetAnimeSession(ctx, info.MalID)
	if cachedSession != nil {
		return cachedSession, nil
	}

	// find the most reliable title to use for searching.
	var title string
	switch {
	case info.TitleJapanese != "":
		title = info.TitleJapanese
	case info.TitleEnglish != "":
		title = info.TitleEnglish
	case info.TitleSynonyms != nil && len(info.TitleSynonyms) > 0:
		for _, t := range info.TitleSynonyms {
			if t == "" {
				continue
			}

			title = t
			break
		}
	default:
		// this quite literally should NEVER happen.
		return nil, errors.New("no title found???")
	}

	search, resp, err := c.Search(ctx, title, 1)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("animepahe api returned status code %d", resp.StatusCode)
	}
	if len(search.Data) <= 0 {
		return nil, ErrNoSearchResultsFound
	}

	var match *SearchResult
	for _, entry := range search.Data {
		// if it's not the same year, ignore
		if entry.Year != info.Year {
			continue
		}

		// if it's not the same season, ignore
		if strings.ToLower(info.Season) != strings.ToLower(entry.Season) {
			continue
		}

		// if it's not the same type, ignore
		if strings.ToLower(info.Type) != strings.ToLower(entry.Type) {
			continue
		}

		// find the most similar title in the array of titles
		var similarity float64
		var closestMatch *SearchResult
		for _, t := range info.Titles {
			titleSimilarty := c.stringCompare(t.Title, entry.Title)

			if titleSimilarty >= similarity {
				similarity = titleSimilarty
				closestMatch = &entry
			}
		}

		if similarity >= 0.9 {
			match = closestMatch
			break
		}
	}

	if match == nil {
		return nil, ErrNoSearchResultsFound
	}

	err = c.Cache.SetAnimeSession(ctx, info.MalID, match.Session)
	if err != nil {
		return nil, err
	}

	return &match.Session, nil
}
