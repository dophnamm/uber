package config

import (
	"fmt"
	"log"

	env "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port    int    `env:"PORT" envDefault:"8082"`
	GinMode string `env:"GIN_MODE" envDefault:"debug"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return cfg, err
}

func InitConfig() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}

	conf, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	return &conf
}
