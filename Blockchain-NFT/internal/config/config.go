package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Duration wraps time.Duration to support YAML unmarshalling from strings like "30s".
type Duration struct {
	time.Duration
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	if value == nil {
		d.Duration = 0
		return nil
	}

	switch value.Tag {
	case "!!str":
		dur, err := time.ParseDuration(value.Value)
		if err != nil {
			return fmt.Errorf("config: invalid duration %q: %w", value.Value, err)
		}
		d.Duration = dur
		return nil
	case "!!int":
		// treat integers as seconds
		dur, err := time.ParseDuration(value.Value + "s")
		if err != nil {
			return fmt.Errorf("config: invalid duration %q: %w", value.Value, err)
		}
		d.Duration = dur
		return nil
	default:
		return fmt.Errorf("config: unsupported duration type %s", value.Tag)
	}
}

// FileConfig mirrors the YAML structure for project configuration.
type FileConfig struct {
	Node struct {
		ID   string `yaml:"id"`
		Type string `yaml:"type"`
		Port int    `yaml:"port"`
	} `yaml:"node"`

	Network struct {
		BootstrapPeers []string `yaml:"bootstrap_peers"`
		Consensus      string   `yaml:"consensus"`
	} `yaml:"network"`

	Supabase struct {
		URL            string   `yaml:"url"`
		AnonKey        string   `yaml:"anon_key"`
		ServiceRoleKey string   `yaml:"service_role_key"`
		PollInterval   Duration `yaml:"poll_interval"`
	} `yaml:"supabase"`

	Storage struct {
		BadgerPath string `yaml:"badger_path"`
		IPFSAPI    string `yaml:"ipfs_api"`
	} `yaml:"storage"`

	API struct {
		EnableWebsocket bool     `yaml:"enable_websocket"`
		RateLimit       int      `yaml:"rate_limit"`
		AllowedOrigins  []string `yaml:"allowed_origins"`
	} `yaml:"api"`
}

// LoadFile parses the configuration YAML at the specified path.
func LoadFile(path string) (*FileConfig, error) {
	if path == "" {
		return nil, errors.New("config: path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file failed: %w", err)
	}

	var cfg FileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: decode yaml failed: %w", err)
	}

	return &cfg, nil
}
