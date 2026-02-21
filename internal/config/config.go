package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	Repo   RepoConfig
	Docs   DocsConfig
}

type ServerConfig struct {
	Port string
}

type RepoConfig struct {
	URL    string
	Branch string
	SSHKey string
}

type DocsConfig struct {
	RefreshInterval time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("APP_PORT", "3000"),
		},
		Repo: RepoConfig{
			URL:    getEnvOrDefault("GIT_REPO_URL", "https://github.com/vacano-house/vacano-ui.git"),
			Branch: getEnvOrDefault("GIT_BRANCH", "master"),
			SSHKey: os.Getenv("GIT_SSH_KEY"),
		},
		Docs: DocsConfig{
			RefreshInterval: parseDuration(getEnvOrDefault("DOCS_REFRESH_INTERVAL", "5m")),
		},
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic("invalid duration format: " + s)
	}
	return d
}
