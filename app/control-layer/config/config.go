package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port      string `env:"PORT" envDefault:"8080"`
	MongoUri  string `env:"MONGO_URI"`
	MongoDb   string `env:"MONGO_DB"`
	RedisUri  string `env:"REDIS_URI"`
	JwtSecret string `env:"JWT_SECRET_CONTROL_PLANE"`
}

func LoadEnv() (Config, error) {
	_ = godotenv.Load()
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
