package config

import (
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	storageType string
	config      struct {
		Host     string
		Port     uint16
		Brand    string
		Blizzard blizzardConfig
		Polar    polarConfig
		Debug    bool `yaml:"-"`
		Orca     orcaConfig
	}

	polarConfig struct {
		Port     uint16
		Password string
	}

	blizzardConfig struct {
		PrivateKey string `yaml:"privateKey"`
		EnableCORS bool   `json:"enableCors"`
		RateLimit  uint32 `yaml:"rateLimit"`
		Database   databaseConfig
		Storage    map[storageType]string
		OAuth      map[string]oauthProvider
		Dragonfly  dragonflyConfig
	}

	orcaConfig struct {
		Token string
		Guild string
	}

	dragonflyConfig struct {
		Host     string
		Port     uint16
		Password string
	}

	oauthProvider struct {
		ClientID     string `yaml:"clientId"`
		ClientSecret string `yaml:"clientSecret"`
	}

	databaseConfig struct {
		Host     string
		Port     uint16
		Username string
		Password string
		Name     string
		Secure   bool
	}
)

var Config config

const (
	Problems    storageType = "problems"
	Posts                   = "posts"
	READMEs                 = "readmes"
	Submissions             = "submissions"
)

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
