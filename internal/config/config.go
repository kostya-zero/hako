package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Address         string `json:"address"`
	SnapshotFile    string `json:"snapshot_file"`
	SnapshotEnabled bool   `json:"snapshot_enabled"`
}

func GetDefaultConfig() Config {
	return Config{
		Address:         ":3000",
		SnapshotFile:    "hako-snapshot.dat",
		SnapshotEnabled: false,
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
