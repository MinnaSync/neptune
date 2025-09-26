package config

import (
	"sync"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port int `env:"PORT" envDefault:"8080"`

	ProviderURLs struct {
		Animepahe string `env:"ANIMEPAHE,required"`
	} `envPrefix:"PROVIDER_URL_"`

	Redis struct {
		Host     string `env:"HOST,required"`
		Port     int    `env:"PORT,required"`
		Password string `env:"PASSWORD" envDefault:"NOPASS"`
		DB       int    `env:"DB" default:"0"`
	} `envPrefix:"REDIS_"`
}

var (
	once sync.Once
	C    Config
)

func setup() {
	if err := env.Parse(&C); err != nil {
		panic(err)
	}
}

func init() {
	once.Do(setup)
}
