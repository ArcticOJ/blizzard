package config

import (
	"blizzard/logger"
	"blizzard/models"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	"os"
)

var Config models.BlizzardConfig

func init() {
	b, e := os.ReadFile("arctic.yml")
	logger.Panic(e, "could not read config file")
	logger.Panic(yaml.Unmarshal(b, &Config), "failed to parse config file")
	Config.Debug = os.Getenv("ENV") == "dev"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if Config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
