package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/minna-sync/neptune/config"
	"github.com/minna-sync/neptune/handlers"
	_ "github.com/minna-sync/neptune/internal/logger"
	uc "github.com/minna-sync/neptune/usecase"
	"github.com/redis/go-redis/v9"
)

func initRedis() (*redis.Client, error) {
	var opts = &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.C.Redis.Host, config.C.Redis.Port),
		Password: config.C.Redis.Password,
		DB:       config.C.Redis.DB,
	}

	client := redis.NewClient(opts)

	return client, nil
}

func main() {
	redis, err := initRedis()
	if err != nil {
		panic(err)
	}

	usecase := uc.NewAnimeLookupUsecase(redis)

	router := chi.NewRouter()
	router.Use(handlers.RequestLog)
	handlers.NewAnimeRouter(router, usecase)

	port := fmt.Sprintf(":%d", config.C.Port)
	panic(http.ListenAndServe(port, router))
}
