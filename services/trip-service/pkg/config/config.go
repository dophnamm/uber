package config

import (
	"log"
	"os"

	env "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port            int    `env:"PORT" envDefault:"8082"`
	GinMode         string `env:"GIN_MODE" envDefault:"debug"`
	MongoPort       string `env:"MONGO_PORT" envDefault:"27017"`
	MongoUri        string `env:"MONGO_URI" envDefault:"mongodb://localhost:27017"`
	MongoDBName     string `env:"MONGO_DB_NAME" envDefault:"ride_sharing"`
	MongoAuthSource string `env:"MONGO_AUTH_SOURCE"`
	MongoUserName   string `env:"MONGO_USER_NAME"`
	MongoPassword   string `env:"MONGO_PASSWORD"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return cfg, err
}

func InitConfig() *Config {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("loading .env: %v", err)
	}

	conf, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	return &conf
}
