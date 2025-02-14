package config

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	// Set up test environment variables
	if err := os.Setenv("SOWESIGN_CODE_ETABLISSEMENT", "test-code"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("SOWESIGN_IDENTIFIANT", "test-id"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("SOWESIGN_PIN", "test-pin"); err != nil {
		t.Fatal(err)
	}

	config, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() failed: %v", err)
	}

	if config.CodeEtablissement == "" {
		t.Error("CodeEtablissement should not be empty")
	}
	if config.Identifiant == "" {
		t.Error("Identifiant should not be empty")
	}
	if config.PIN == "" {
		t.Error("PIN should not be empty")
	}
}

func TestNewConfig_MissingValues(t *testing.T) {
	// Clear environment variables
	for _, env := range []string{
		"SOWESIGN_CODE_ETABLISSEMENT",
		"SOWESIGN_IDENTIFIANT",
		"SOWESIGN_PIN",
	} {
		if err := os.Unsetenv(env); err != nil {
			t.Fatal(err)
		}
	}

	_, err := NewConfig()
	if err == nil {
		t.Error("Expected error with missing environment variables")
	}
}

func TestNewTestConfig(t *testing.T) {
	config := NewTestConfig()

	if config.CodeEtablissement != "test-code" {
		t.Errorf("Expected test-code, got %s", config.CodeEtablissement)
	}
	if config.Identifiant != "test-id" {
		t.Errorf("Expected test-id, got %s", config.Identifiant)
	}
	if config.PIN != "test-pin" {
		t.Errorf("Expected test-pin, got %s", config.PIN)
	}
}
