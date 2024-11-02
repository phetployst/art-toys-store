package config

import "testing"

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
			t.Errorf("expected %q but got %q", got, want)
		}
	})

	t.Run("get default value given key does not exist", func(t *testing.T) {
		got := configProvider.GetStringEnv("NOT_EXIST", "world")
		want := "world"

		if got != want {
			t.Errorf("expected %q but got %q", got, want)
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
			t.Errorf("expected %d but got %d", got, want)
		}
	})

	t.Run("get default value given key does not exist", func(t *testing.T) {
		envGetter := StubEnvGetter{}
		configProvider := ConfigProvider{Getter: envGetter}
		got := configProvider.GetIntEnv("NOT_EXIST", 10)
		want := 10

		if got != want {
			t.Errorf("expected %d but got %d", got, want)
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
			t.Errorf("expected %d but got %d", got, want)
		}
	})
}

func TestGetConfig(t *testing.T) {
	t.Run("get server given keys exist", func(t *testing.T) {
		envGetter := StubEnvGetter{
			"ENVIRONMENT":          "local",
			"SERVICE_NAME":         "auth",
			"HOSTNAME":             "localhost",
			"PORT":                 "5000",
			"DB_CONNECTION_STRING": "db://localhost:5432",
		}
		configProvider := ConfigProvider{Getter: envGetter}
		config := configProvider.GetConfig()

		got := config
		want := Config{
			Environment: "local",
			Server: Server{
				ServiceName:        "auth",
				Hostname:           "localhost",
				Port:               5000,
				DBConnectionString: "db://localhost:5432",
			},
		}

		if got != want {
			t.Errorf("expected %v but got %v", got, want)
		}
	})

	t.Run("get server given keys do not exist", func(t *testing.T) {
		envGetter := StubEnvGetter{}
		configProvider := ConfigProvider{Getter: envGetter}
		config := configProvider.GetConfig()

		got := config
		want := Config{
			Environment: "local",
			Server: Server{
				ServiceName:        "",
				Hostname:           "localhost",
				Port:               0,
				DBConnectionString: "",
			},
		}

		if got != want {
			t.Errorf("expected %v but got %v", got, want)
		}
	})
}
