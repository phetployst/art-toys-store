package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	EnvGetter interface {
		Getenv(key string) string
	}

	OsEnvGetter struct{}

	ConfigProvider struct {
		Getter EnvGetter
	}

	Config struct {
		Environment string
		Server      Server
		Jwt         Jwt
	}

	Server struct {
		ServiceName        string
		Hostname           string
		Port               int
		DBConnectionString string
	}

	Jwt struct {
		AccessTokenSecret  string
		RefreshTokenSecret string
	}
)

func (o *OsEnvGetter) Getenv(key string) string {
	return os.Getenv(key)
}

func (c *ConfigProvider) GetStringEnv(key string, defaultValue string) string {
	value := c.Getter.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *ConfigProvider) GetIntEnv(key string, defaultValue int) int {
	value := c.Getter.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func (c *ConfigProvider) GetRequiredEnv(key string) (string, error) {
	value := c.Getter.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return value, nil
}

func LoadEnvFile(path string) {
	if err := godotenv.Load(path); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func (c *ConfigProvider) GetConfig() (Config, error) {

	accessTokenSecret, err := c.GetRequiredEnv("JWT_ACCESS_SECRET")
	if err != nil {
		return Config{}, fmt.Errorf("failed to load JWT_ACCESS_SECRET: %w", err)
	}

	refreshTokenSecret, err := c.GetRequiredEnv("JWT_REFRESH_SECRET")
	if err != nil {
		return Config{}, fmt.Errorf("failed to load JWT_REFRESH_SECRET: %w", err)
	}

	return Config{
		Environment: c.GetStringEnv("ENVIRONMENT", "local"),
		Server: Server{
			ServiceName:        c.GetStringEnv("SERVICE_NAME", "account"),
			Hostname:           c.GetStringEnv("HOSTNAME", "localhost"),
			Port:               c.GetIntEnv("PORT", 1323),
			DBConnectionString: c.GetStringEnv("DB_CONNECTION_STRING", ""),
		},
		Jwt: Jwt{
			AccessTokenSecret:  accessTokenSecret,
			RefreshTokenSecret: refreshTokenSecret,
		},
	}, nil
}
