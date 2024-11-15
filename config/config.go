package config

import (
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

func LoadEnvFile(path string) {
	if err := godotenv.Load(path); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func (c *ConfigProvider) GetConfig() Config {
	return Config{
		Environment: c.GetStringEnv("ENVIRONMENT", "local"),
		Server: Server{
			ServiceName:        c.GetStringEnv("SERVICE_NAME", ""),
			Hostname:           c.GetStringEnv("HOSTNAME", "localhost"),
			Port:               c.GetIntEnv("PORT", 0),
			DBConnectionString: c.GetStringEnv("DB_CONNECTION_STRING", ""),
		},
		Jwt: Jwt{
			AccessTokenSecret:  c.GetStringEnv("JWT_ACCESS_SECRET", ""),
			RefreshTokenSecret: c.GetStringEnv("JWT_REFRESH_SECRET", ""),
		},
	}
}
