package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CodeEtablissement string `json:"codeEtablissement"`
	Identifiant       string `json:"identifiant"`
	PIN               string `json:"PIN"`
}

// NewConfig creates a new Config instance from environment variables
func NewConfig() (Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Only return error if file exists but couldn't be loaded
		if !os.IsNotExist(err) {
			return Config{}, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	cfg := Config{
		CodeEtablissement: os.Getenv("SOWESIGN_CODE_ETABLISSEMENT"),
		Identifiant:       os.Getenv("SOWESIGN_IDENTIFIANT"),
		PIN:               os.Getenv("SOWESIGN_PIN"),
	}

	// Validate required fields
	if cfg.CodeEtablissement == "" {
		return Config{}, fmt.Errorf("SOWESIGN_CODE_ETABLISSEMENT is required")
	}
	if cfg.Identifiant == "" {
		return Config{}, fmt.Errorf("SOWESIGN_IDENTIFIANT is required")
	}
	if cfg.PIN == "" {
		return Config{}, fmt.Errorf("SOWESIGN_PIN is required")
	}

	return cfg, nil
}

// NewTestConfig creates a Config instance for testing
func NewTestConfig() Config {
	return Config{
		CodeEtablissement: "test-code",
		Identifiant:       "test-id",
		PIN:               "test-pin",
	}
}
