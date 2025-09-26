package jikan

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Anime service

type AnimeImageURLs struct {
	ImageURL      string `json:"image_url"`
	SmallImageURL string `json:"small_image_url"`
	LargeImageURL string `json:"large_image_url"`
}

type AnimeImages struct {
	JPG  AnimeImageURLs `json:"jpg"`
	WebP AnimeImageURLs `json:"webp"`
}

type AnimeTrailer struct {
	YouTubeID string `json:"youtube_id"`
	URL       string `json:"url"`
	EmbedURL  string `json:"embed_url"`
}

type AnimeTitle struct {
	Type  string `json:"type"`
	Title string `json:"title"`
}

type AnimeDate struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type AnimeAiredProperties struct {
	From   AnimeDate `json:"from"`
	To     AnimeDate `json:"to"`
	String string    `json:"string"`
}

type AnimeAired struct {
	From  string               `json:"from"`
	To    string               `json:"to"`
	Props AnimeAiredProperties `json:"props"`
}

type AnimeBroadcast struct {
	Day      string `json:"day"`
	Time     string `json:"time"`
	Timezone string `json:"timezone"`
	String   string `json:"string"`
}

type AnimeMALRelation struct {
	MalID int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type AnimeInfoBase struct {
	MalID          int                `json:"mal_id"`
	URL            string             `json:"url"`
	Images         AnimeImages        `json:"images"`
	Trailer        AnimeTrailer       `json:"trailer"`
	Approved       bool               `json:"approved"`
	Titles         []AnimeTitle       `json:"titles"`
	Title          string             `json:"title"`
	TitleEnglish   string             `json:"title_english"`
	TitleJapanese  string             `json:"title_japanese"`
	TitleSynonyms  []string           `json:"title_synonyms"`
	Type           string             `json:"type"`
	Source         string             `json:"source"`
	Episodes       int                `json:"episodes"`
	Status         string             `json:"status"`
	Airing         bool               `json:"airing"`
	Aired          AnimeAired         `json:"aired"`
	Duration       string             `json:"duration"`
	Rating         string             `json:"rating"`
	Score          float64            `json:"score"`
	ScoreBy        int                `json:"scored_by"`
	Rank           int                `json:"rank"`
	Popularity     int                `json:"popularity"`
	Members        int                `json:"members"`
	Favorites      int                `json:"favorites"`
	Synopsis       string             `json:"synopsis"`
	Background     string             `json:"background"`
	Season         string             `json:"season"`
	Year           int                `json:"year"`
	Broadcast      AnimeBroadcast     `json:"broadcast"`
	Producers      []AnimeMALRelation `json:"producers"`
	Licensors      []AnimeMALRelation `json:"licensors"`
	Studios        []AnimeMALRelation `json:"studios"`
	Genres         []AnimeMALRelation `json:"genres"`
	ExplicitGenres []AnimeMALRelation `json:"explicit_genres"`
	Themes         []AnimeMALRelation `json:"themes"`
	Demographics   []AnimeMALRelation `json:"demographics"`
}

func (s *Anime) GetAnimeFullById() {
}

func (s *Anime) GetAnimeById(ctx context.Context, id int) (*Result[*AnimeInfoBase], *http.Response, error) {
	u := fmt.Sprintf("/v4/anime/%d", id)

	req, err := s.client.NewGetRequest(u)
	if err != nil {
		return nil, nil, err
	}

	info := new(Result[*AnimeInfoBase])
	resp, err := s.client.Do(ctx, req, info)
	if err != nil {
		return nil, resp, err
	}

	return info, resp, nil
}

type AnimeEpisode struct {
	MalID         int    `json:"mal_id"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	TitleJapanese string `json:"title_japanese"`
	TitleEnglish  string `json:"title_english"`
	Duration      int    `json:"duration"`
	Aired         string `json:"aired"`
	Filler        bool   `json:"filler"`
	Recap         bool   `json:"recap"`
	Synopsis      string `json:"synopsis"`
	ForumURL      string `json:"forum_url"`
}

func (s *Anime) GetAnimeEpisodes(ctx context.Context, id int, page int) (*PaginatedResults[[]AnimeEpisode], *http.Response, error) {
	u := fmt.Sprintf("/v4/anime/%d/episodes?page=%d", id, page)

	req, err := s.client.NewGetRequest(u)
	if err != nil {
		return nil, nil, err
	}

	results := new(PaginatedResults[[]AnimeEpisode])
	resp, err := s.client.Do(ctx, req, results)
	if err != nil {
		return nil, resp, err
	}

	return results, resp, nil
}

func (s *Anime) GetAnimeSearch(ctx context.Context, query url.Values) (*PaginatedResults[[]AnimeInfoBase], *http.Response, error) {
	u := fmt.Sprintf("/v4/anime?%s", query.Encode())

	req, err := s.client.NewGetRequest(u)
	if err != nil {
		return nil, nil, err
	}

	results := new(PaginatedResults[[]AnimeInfoBase])
	resp, err := s.client.Do(ctx, req, results)
	if err != nil {
		return nil, nil, err
	}

	return results, resp, nil
}
