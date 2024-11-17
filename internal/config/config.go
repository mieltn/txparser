package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Mode string
	App  struct {
		Port            string `json:"port"`
		PollIntervalSec int    `json:"poll_interval_sec"`
		PollWorkers     int    `json:"poll_workers"`
		StartBlock      string `json:"start_block"`
	} `json:"app"`
	Eth struct {
		Url     string `json:"url"`
		Retry   int    `json:"retry"`
		RetryIn int    `json:"retry_in"`
		Timeout int    `json:"timeout"`
	} `json:"eth"`
}

func Load(cfg *Config) error {
	file, err := os.ReadFile(fmt.Sprintf("internal/config/%s.json", cfg.Mode))
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return err
	}

	return nil
}
