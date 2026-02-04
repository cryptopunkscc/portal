package config

import (
	"fmt"
	"testing"
)

// Sample Config struct with nested structures.
type Config struct {
	Port     int
	Hostname string
	Timeout  float64
	Debug    bool
	Nested   NestedConfig
}

type NestedConfig struct {
	Level    int
	Msg      string
	Internal InnerConfig
}

type InnerConfig struct {
	Flag bool
}

// TestMerge tests the recursive merging of a user config with a default config.
func TestMerge(t *testing.T) {
	// Default config.
	defaultConfig := Config{
		Port:     8080,
		Hostname: "localhost",
		Timeout:  30.0,
		Debug:    true,
		Nested: NestedConfig{
			Level: 1,
			Msg:   "default message",
			Internal: InnerConfig{
				Flag: true,
			},
		},
	}

	// Sample user config that sets only some fields.
	userConfig := Config{
		Port: 3000,
		// Hostname, Timeout, Debug are zero values.
		Nested: NestedConfig{
			// Level is zero, Msg is set.
			Msg:      "custom message",
			Internal: InnerConfig{
				// Flag is zero (false)
			},
		},
	}

	Merge(&userConfig, &defaultConfig)

	expected := Config{
		Port:     3000,        // From userConfig.
		Hostname: "localhost", // Merged from defaultConfig.
		Timeout:  30.0,        // Merged from defaultConfig.
		Debug:    true,        // Merged from defaultConfig.
		Nested: NestedConfig{
			Level: 1,                // Merged from defaultConfig.
			Msg:   "custom message", // From userConfig.
			Internal: InnerConfig{
				Flag: true, // Merged from defaultConfig.
			},
		},
	}

	if userConfig != expected {
		t.Errorf("Merged config is not as expected.\nGot: %+v\nExpected: %+v", userConfig, expected)
	} else {
		t.Logf("Merged Config: %+v", userConfig)
	}

	// For visual feedback when running "go test -v", print the merged config.
	fmt.Printf("Merged Config: %+v\n", userConfig)
}
