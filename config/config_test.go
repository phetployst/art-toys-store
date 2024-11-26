package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StubEnvGetter map[string]string

func (s StubEnvGetter) Getenv(key string) string {
	return s[key]
}

type MockEnvLoader struct {
	LoadenvFunc func(path string) error
}

func (m *MockEnvLoader) Loadenv(path string) error {
	if m.LoadenvFunc != nil {
		return m.LoadenvFunc(path)
	}
	return nil
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

func TestLoadEnvFile(t *testing.T) {
	t.Run("successfully load env file", func(t *testing.T) {
		mockLoader := &MockEnvLoader{
			LoadenvFunc: func(path string) error {
				return nil
			},
		}
		configProvider := ConfigProvider{Loader: mockLoader}

		err := configProvider.LoadEnvFile(".env")
		assert.NoError(t, err, "expected no error when env file is loaded successfully")
	})

	t.Run("return error when env file fails to load", func(t *testing.T) {
		mockLoader := &MockEnvLoader{
			LoadenvFunc: func(path string) error {
				return errors.New("failed to load .env file")
			},
		}
		configProvider := ConfigProvider{Loader: mockLoader}

		err := configProvider.LoadEnvFile(".env")
		assert.Error(t, err, "expected an error when env file fails to load")
		assert.Contains(t, err.Error(), "failed to load .env file")
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

	t.Run("get default value when keys do not exist", func(t *testing.T) {
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
