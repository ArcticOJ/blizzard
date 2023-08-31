package config

import (
	"blizzard/blizzard/logger"
	"blizzard/blizzard/models"
	"gopkg.in/yaml.v3"
	"os"
)

var Config models.BlizzardConfig

func init() {
	b, e := os.ReadFile("blizzard.yml")
	logger.Panic(e, "could not read config file")
	logger.Panic(yaml.Unmarshal(b, &Config), "failed to parse config file")
}
