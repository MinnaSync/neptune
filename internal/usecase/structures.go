package usecase

import "github.com/minna-sync/neptune/pkg/jikan"

type AnimeTitle struct {
	Native  *string `json:"native"`
	Romaji  *string `json:"romaji"`
	English *string `json:"english"`
}

type AnimeSearchEntry struct {
	MALId int `json:"mal_id"`

	Poster string     `json:"poster"`
	Title  AnimeTitle `json:"title"`
	Type   string     `json:"type"`
	Year   int        `json:"year"`

	AiredEpisodes int `json:"aired_episodes"`
	TotalEpisodes int `json:"episodes"`

	IsNSFW bool `json:"is_nsfw"`
}

type AnimeInfo struct {
	Info *jikan.AnimeInfoBase `json:"info"`
}

type AnimeEpisode struct {
	ID       any     `json:"id"` // could be string or int
	Episode  int     `json:"episode"`
	Title    string  `json:"title"`
	Synopsis *string `json:"synopsis"`
	Snapshot string  `json:"snapshot"`
}

type AnimeEpisodeSubtitles struct {
	Language string `json:"language"`
	URL      string `json:"url"`
}

type AnimeEpisodeLinks struct {
	URL        string                  `json:"url"`
	Resolution string                  `json:"resolution"`
	Language   string                  `json:"language"`
	Subtitles  []AnimeEpisodeSubtitles `json:"subtitles"`
}
