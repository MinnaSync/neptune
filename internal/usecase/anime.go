package usecase

import "context"

type AnimeLookupUsecase interface {
	SearchAnime(ctx context.Context, term string) ([]AnimeSearchEntry, error)
	GetAnime(ctx context.Context, id int) (*AnimeInfo, error)
	GetEpisodes(ctx context.Context, id int, provider string, page int) ([]AnimeEpisode, error)
	GetStreamingLinks(ctx context.Context, id int, provider string, episode string) ([]AnimeEpisodeLinks, error)
}
