package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port             int    `json:"port"`
	SnapshotFile     string `json:"snapshot_file"`
	SnapshotsEnabled bool   `json:"snapshots_enabled"`
}

func GetDefaultConfig() Config {
	return Config{
		Port:             7000,
		SnapshotFile:     "hako-snapshot.dat",
		SnapshotsEnabled: false,
	}
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
