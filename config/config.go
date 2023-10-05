package config

import (
	"blizzard/logger"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	BlizzardConfig struct {
		Address    `yaml:",inline"`
		PrivateKey string `yaml:"privateKey"`
		Debug      bool   `yaml:"-"`
		EnableCORS bool   `json:"enableCors"`
		RateLimit  uint32 `yaml:"rateLimit"`
		Database   DatabaseConfig
		Storage    StorageConfig
		Discord    *DiscordConfig
		OAuth      map[string]OAuthProvider
		Dragonfly  Address
		RabbitMQ   RabbitMQConfig
	}

	DiscordConfig struct {
		Token string
		Guild string
	}

	RabbitMQConfig struct {
		Username    string
		Password    string
		Address     `yaml:",inline"`
		ManagerPort uint16 `yaml:"managerPort"`
		StreamPort  uint16 `yaml:"streamPort"`
		VHost       string
	}

	Address struct {
		Host string
		Port uint16
	}

	OAuthProvider struct {
		ClientID     string `yaml:"clientId"`
		ClientSecret string `yaml:"clientSecret"`
	}

	StorageConfig struct {
		Problems    string
		Posts       string
		READMEs     string
		Submissions string
	}

	DatabaseConfig struct {
		Address  `yaml:",inline"`
		Username string
		Password string
		Name     string
		Secure   bool
	}
)

var Config BlizzardConfig

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
