package config

import (
	"blizzard/blizzard/logger"
	"blizzard/blizzard/models"
	"github.com/spf13/viper"
)

var Config *models.BlizzardConfig

func readConfig() *models.BlizzardConfig {
	// TODO: Command line arguments, env config and config file
	var conf models.BlizzardConfig
	v := viper.New()
	v.SetConfigName("blizzard")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
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
