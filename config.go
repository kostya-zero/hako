package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port             int
	SnapshotFile     string
	SnapshotsEnabled bool `json:"snapshots_enabled"`
}

func GetDefaultConfig() Config {
	return Config{
		Port:             7000,
		SnapshotFile:     "hako-snapshot.dat",
		SnapshotsEnabled: true,
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
