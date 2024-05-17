package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type DbConfig struct {
	Path string `envconfig:"DB_PATH" required:"true"`
}

type JwtConfig struct {
	PrivateKey string `envconfig:"JWT_PRIVATE_KEY" required:"true"`
}

type Config struct {
	Port int `envconfig:"PORT" required:"true"`
	DB   DbConfig
	Jwt  JwtConfig
}

func GetConfigFromEnv(path string) Config {
	var config Config
	godotenv.Load(".env")
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}

	return config
}
