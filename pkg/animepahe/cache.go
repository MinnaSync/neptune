package animepahe

import (
	"context"
	"fmt"
	"time"

	"github.com/minna-sync/neptune/pkg/redisx"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type AnimepaheCache struct {
	sf    singleflight.Group
	redis *redis.Client
}

func NewAnimepaheCache(redis *redis.Client) *AnimepaheCache {
	return &AnimepaheCache{redis: redis}
}

func (c *AnimepaheCache) SetAnimeSession(ctx context.Context, animeId int, session string) error {
	key := fmt.Sprintf("animepahe:%d:session", animeId)

	status := c.redis.Set(ctx, key, session, 0)

	if err := status.Err(); err != nil {
		return err
	}

	return nil
}

func (c *AnimepaheCache) GetAnimeSession(ctx context.Context, animeId int) (*string, error) {
	key := fmt.Sprintf("animepahe:%d:session", animeId)

	result, err, _ := c.sf.Do(key, func() (any, error) {
		result := c.redis.Get(ctx, key)

		if err := result.Err(); err != nil {
			return nil, err
		}

		val := result.Val()

		return &val, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*string), err
}

func (c *AnimepaheCache) SetEpisodes(ctx context.Context, animeId int, episodes *[]EpisodeResult, page int) error {
	key := fmt.Sprintf("animepahe:%d:episodes_%d", animeId, page)

	pipeline := c.redis.Pipeline()
	pipeline.JSONSet(ctx, key, "$", episodes)
	pipeline.Expire(ctx, key, time.Hour*24)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *AnimepaheCache) GetEpisodes(ctx context.Context, animeId int, page int) (*[]EpisodeResult, error) {
	key := fmt.Sprintf("animepahe:%d:episodes_%d", animeId, page)

	result, err, _ := c.sf.Do(key, func() (any, error) {
		var episodes []EpisodeResult
		err := redisx.JSONUnwrap(ctx, c.redis, key, "$", &episodes)

		if err != nil {
			if err == redis.Nil {
				return nil, nil
			}

			return nil, err
		}

		return &episodes, nil
	})

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result.(*[]EpisodeResult), nil
}

// func (c *AnimepaheCache) SetStreamingLink(ctx context.Context, animeId int, episodeId string, link EpisodeStreamingLink) error {
// 	setKey := fmt.Sprintf("animepahe:%d:streaming:%s", animeId, episodeId)
// 	jsonKey := fmt.Sprintf("animepahe:%d:streaming:%s:%s_%s", animeId, episodeId, link.Resolution, link.Language)

// 	pipeline := c.redis.Pipeline()
// 	pipeline.SAdd(ctx, setKey, jsonKey, time.Hour*12)
// 	pipeline.JSONSet(ctx, jsonKey, "$", link)

// 	_, err := pipeline.Exec(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (c *AnimepaheCache) SetStreamingLinks(ctx context.Context, animeId int, episodeId string, links []*EpisodeStreamingLink) error {
	key := fmt.Sprintf("animepahe:%d:streaming:%s", animeId, episodeId)

	pipeline := c.redis.Pipeline()
	pipeline.JSONSet(ctx, key, "$", links)
	pipeline.Expire(ctx, key, time.Hour*24)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *AnimepaheCache) GetStreamingLinks(ctx context.Context, animeId int, episodeId string) ([]*EpisodeStreamingLink, error) {
	key := fmt.Sprintf("animepahe:%d:streaming:%s", animeId, episodeId)

	result, err, _ := c.sf.Do(key, func() (any, error) {
		var links []*EpisodeStreamingLink

		err := redisx.JSONUnwrap(ctx, c.redis, key, "$", &links)

		if err != nil {
			return nil, err
		}

		return links, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*EpisodeStreamingLink), nil
}
