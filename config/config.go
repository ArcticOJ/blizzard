package config

import (
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type (
	storageType string
	config      struct {
		Host     string         `yaml:"host"`
		Port     uint16         `yaml:"port"`
		Brand    string         `yaml:"brand"`
		Blizzard blizzardConfig `yaml:"blizzard"`
		Polar    polarConfig    `yaml:"polar"`
		Debug    bool           `yaml:"-"`
		Orca     orcaConfig     `yaml:"orca"`
	}

	polarConfig struct {
		Port     uint16 `yaml:"port"`
		Secret   string `yaml:"secret"`
		CertFile string `yaml:"certFile"`
		KeyFile  string `yaml:"keyFile"`
	}

	blizzardConfig struct {
		Secret     string                   `yaml:"secret"`
		EnableCORS bool                     `json:"enableCors"`
		RateLimit  uint32                   `yaml:"rateLimit"`
		Database   databaseConfig           `yaml:"database"`
		Storage    map[storageType]string   `yaml:"storage"`
		OAuth      map[string]oauthProvider `yaml:"oauth"`
		Dragonfly  dragonflyConfig          `yaml:"dragonfly"`
	}

	orcaConfig struct {
		Token string `yaml:"token"`
		Guild string `yaml:"guild"`
	}

	dragonflyConfig struct {
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
		Password string `yaml:"password"`
	}

	oauthProvider struct {
		ClientID     string `yaml:"clientId"`
		ClientSecret string `yaml:"clientSecret"`
	}

	databaseConfig struct {
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Secure   bool   `yaml:"secure"`
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
	confPath := strings.TrimSpace(os.Getenv("ARCTIC_CONFIG_PATH"))
	if confPath == "" {
		confPath = "arctic.yml"
	}
	b, e := os.ReadFile(confPath)
	logger.Panic(e, "could not read config file")
	logger.Panic(yaml.Unmarshal(b, &Config), "failed to parse config file")
	Config.Debug = os.Getenv("ENV") == "dev"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if Config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
