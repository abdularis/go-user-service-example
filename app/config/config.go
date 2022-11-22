package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	JWTSecretAccessToken  string `envconfig:"JWT_SECRET_ACCESS_TOKEN" default:"jwt-secret-for-access-token"`
	JWTSecretRefreshToken string `envconfig:"JWT_SECRET_REFRESH_TOKEN" default:"jwt-secret-for-refresh-token"`
	DbHost                string `envconfig:"DB_HOST" default:"localhost"`
	DbPort                uint16 `envconfig:"DB_PORT" default:"3306"`
	DbUser                string `envconfig:"DB_USER" default:"root"`
	DbPassword            string `envconfig:"DB_PASSWORD" default:"root"`
	DbName                string `envconfig:"DB_NAME" default:"go_user_service"`
	HostPort              uint16 `envconfig:"HOST_PORT" default:"8080"`
}

func Get() Config {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatal(err)
	}
	return *cfg
}
