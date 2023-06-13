package config

import (
	"blizzard/blizzard/logger"
	"blizzard/blizzard/models"
	"blizzard/blizzard/utils"
	"github.com/spf13/viper"
)

var Config *models.BlizzardConfig

// TODO: finalize defaultConfig
var defaultConfig = map[string]interface{}{
	"host": "0.0.0.0",
	"port": 2999,
	// TODO: use a machine-bound key as privateKey instead of cryptographically random key
	"privateKey": utils.Rand(16, ""),
	"debug":      false,
	"enableCors": true,
	"rateLimit":  1000,
	"judges":     map[string]models.Judge{},
	"oauth":      map[string]models.OAuthProvider{},
	"storage":    models.StorageConfig{},
	"database": models.DatabaseConfig{
		Address:  "localhost:5432",
		Name:     "postgres",
		Username: "postgres",
		Password: "postgres",
		Secure:   false,
	},
}

func readConfig() *models.BlizzardConfig {
	// TODO: Command line arguments, env config and config file
	var conf models.BlizzardConfig
	v := viper.New()
	v.SetConfigName("blizzard")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	for key, val := range defaultConfig {
		v.SetDefault(key, val)
	}
	if e := v.ReadInConfig(); e != nil {
		logger.Logger.Err(e).Msg("config_reader")
	}
	v.SetEnvPrefix("BLIZZARD")
	v.AutomaticEnv()
	err := v.Unmarshal(&conf)
	if err != nil {
		logger.Logger.Err(err).Msg("config_unmarshaler")
	}
	return &conf
}

func Load() {
	Config = readConfig()
}
