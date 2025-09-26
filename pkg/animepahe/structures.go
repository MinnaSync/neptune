package animepahe

type PaginatedResponse[T any] struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	From        int `json:"from"`
	To          int `json:"to"`
	Data        T   `json:"data"`
}

type SearchResult struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Type     string  `json:"type"`
	Episodes int     `json:"episodes"`
	Status   string  `json:"status"`
	Season   string  `json:"season"`
	Year     int     `json:"year"`
	Score    float64 `json:"score"`
	Poster   string  `json:"poster"`
	Session  string  `json:"session"`
}

type SearchResponse PaginatedResponse[[]SearchResult]

type EpisodeResult struct {
	ID        int    `json:"id"`
	AnimeID   int    `json:"anime_id"`
	Episode   int    `json:"episode"`
	Episode2  int    `json:"episode2"`
	Edition   string `json:"edition"`
	Title     string `json:"title"`
	Snapshot  string `json:"snapshot"`
	Disc      string `json:"disc"`
	Audio     string `json:"audio"`
	Duration  string `json:"duration"`
	Session   string `json:"session"`
	Filter    int    `json:"filter"`
	CreatedAt string `json:"created_at"`
}

type EpisodeResponse PaginatedResponse[[]EpisodeResult]

type EpisodeStreamingLink struct {
	URL        string `json:"url"`
	Resolution string `json:"resolution"`
	Language   string `json:"language"`
}
