package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type StubEnvGetter map[string]string

func (s StubEnvGetter) Getenv(key string) string {
	return s[key]
}

func TestGetStringEnv(t *testing.T) {
	envGetter := StubEnvGetter{
		"STRING_VAR": "hello",
	}
	configProvider := ConfigProvider{Getter: envGetter}
	t.Run("get value given key exists", func(t *testing.T) {
		got := configProvider.GetStringEnv("STRING_VAR", "")
		want := "hello"

		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}
	})

	t.Run("get default value given key does not exist", func(t *testing.T) {
		got := configProvider.GetStringEnv("NOT_EXIST", "world")
		want := "world"

		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}
	})
}

func TestGetIntEnv(t *testing.T) {
	t.Run("get value given key exists", func(t *testing.T) {
		envGetter := StubEnvGetter{
			"INT_VAR": "42",
		}
		configProvider := ConfigProvider{Getter: envGetter}
		got := configProvider.GetIntEnv("INT_VAR", 0)
		want := 42

		if got != want {
			t.Errorf("expected %d but got %d", want, got)
		}
	})

	t.Run("get default value given key does not exist", func(t *testing.T) {
		envGetter := StubEnvGetter{}
		configProvider := ConfigProvider{Getter: envGetter}
		got := configProvider.GetIntEnv("NOT_EXIST", 10)
		want := 10

		if got != want {
			t.Errorf("expected %d but got %d", want, got)
		}
	})

	t.Run("get default value given value cannot be converted to int", func(t *testing.T) {
		envGetter := StubEnvGetter{
			"INT_VAR": "42.5",
		}
		configProvider := ConfigProvider{Getter: envGetter}
		got := configProvider.GetIntEnv("INT_VAR", 10)
		want := 10

		if got != want {
			t.Errorf("expected %d but got %d", want, got)
		}
	})
}

func TestGetRequiredEnv(t *testing.T) {
	t.Run("get value given key exists", func(t *testing.T) {
		envGetter := StubEnvGetter{
			"REQUIRED_VAR": "SECRET",
		}
		configProvider := ConfigProvider{Getter: envGetter}
		got, _ := configProvider.GetRequiredEnv("REQUIRED_VAR")
		want := "SECRET"

		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("get value given error when key do not exists", func(t *testing.T) {
		envGetter := StubEnvGetter{}
		configProvider := ConfigProvider{Getter: envGetter}
		_, err := configProvider.GetRequiredEnv("REQUIRED_VAR")

		assert.Error(t, err)
	})
}

func TestGetConfig(t *testing.T) {
	t.Run("get env given keys exist", func(t *testing.T) {
		envGetter := StubEnvGetter{
			"ENVIRONMENT":          "local",
			"SERVICE_NAME":         "auth",
			"HOSTNAME":             "localhost",
			"PORT":                 "5000",
			"DB_CONNECTION_STRING": "db://localhost:5432",
			"JWT_ACCESS_SECRET":    "access-secret",
			"JWT_REFRESH_SECRET":   "refresh-secret",
		}
		configProvider := ConfigProvider{Getter: envGetter}
		config, err := configProvider.GetConfig()

		got := config
		want := Config{
			Environment: "local",
			Server: Server{
				ServiceName:        "auth",
				Hostname:           "localhost",
				Port:               5000,
				DBConnectionString: "db://localhost:5432",
			},
			Jwt: Jwt{
				AccessTokenSecret:  "access-secret",
				RefreshTokenSecret: "refresh-secret",
			},
		}

		assert.NoError(t, err)

		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("get default value server when keys do not exist", func(t *testing.T) {
		envGetter := StubEnvGetter{
			"JWT_ACCESS_SECRET":  "access-secret",
			"JWT_REFRESH_SECRET": "refresh-secret",
		}
		configProvider := ConfigProvider{Getter: envGetter}
		config, err := configProvider.GetConfig()

		got := config
		want := Config{
			Environment: "local",
			Server: Server{
				ServiceName:        "account",
				Hostname:           "localhost",
				Port:               1323,
				DBConnectionString: "",
			},
			Jwt: Jwt{
				AccessTokenSecret:  "access-secret",
				RefreshTokenSecret: "refresh-secret",
			},
		}

		assert.NoError(t, err)

		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("get error given JWT secret do not exist", func(t *testing.T) {
		envGetter := StubEnvGetter{}
		configProvider := ConfigProvider{Getter: envGetter}
		_, err := configProvider.GetConfig()

		assert.Error(t, err)
	})
}
