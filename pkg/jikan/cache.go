package jikan

import (
	"context"
	"fmt"
	"time"

	"github.com/minna-sync/neptune/pkg/redisx"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type JikanCache struct {
	sf    singleflight.Group
	redis *redis.Client
}

func NewJikanCache(redis *redis.Client) *JikanCache {
	return &JikanCache{redis: redis}
}

func (c *JikanCache) SetAnimeInfo(ctx context.Context, id int, anime *AnimeInfoBase) error {
	key := fmt.Sprintf("jikan:anime:%d:info", id)

	pipeline := c.redis.Pipeline()
	pipeline.JSONSet(ctx, key, "$", anime)
	pipeline.Expire(ctx, key, time.Hour*24)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *JikanCache) GetAnimeInfo(ctx context.Context, id int) (*AnimeInfoBase, error) {
	key := fmt.Sprintf("jikan:anime:%d:info", id)

	result, err, _ := c.sf.Do(key, func() (any, error) {
		var info AnimeInfoBase
		err := redisx.JSONUnwrap(ctx, c.redis, key, "$", &info)

		if err != nil {
			return nil, err
		}

		return &info, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*AnimeInfoBase), nil
}
