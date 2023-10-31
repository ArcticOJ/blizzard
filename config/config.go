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
		address    `yaml:",inline"`
		Brand      string
		PrivateKey string `yaml:"privateKey"`
		Debug      bool   `yaml:"-"`
		EnableCORS bool   `json:"enableCors"`
		RateLimit  uint32 `yaml:"rateLimit"`
		Database   databaseConfig
		Storage    map[storageType]string
		Discord    *discordConfig
		OAuth      map[string]oauthProvider
		Dragonfly  address
		RabbitMQ   rabbitMqConfig
	}

	discordConfig struct {
		Token string
		Guild string
	}

	rabbitMqConfig struct {
		Username    string
		Password    string
		address     `yaml:",inline"`
		ManagerPort uint16 `yaml:"managerPort"`
		StreamPort  uint16 `yaml:"streamPort"`
		VHost       string
	}

	address struct {
		Host string
		Port uint16
	}

	oauthProvider struct {
		ClientID     string `yaml:"clientId"`
		ClientSecret string `yaml:"clientSecret"`
	}

	databaseConfig struct {
		address  `yaml:",inline"`
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
