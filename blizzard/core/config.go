package core

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Host       string            `yaml:"host"`
	PrivateKey string            `yaml:"privateKey"`
	Port       uint16            `yaml:"port"`
	Debug      bool              `yaml:"debug"`
	EnableCORS bool              `yaml:"cors"`
	RateLimit  float64           `yaml:"rateLimit"`
	Database   DatabaseConfig    `yaml:"database"`
	Judges     map[string]string `yaml:"judges"`
}

type DatabaseConfig struct {
	Address      string `yaml:"address"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"name"`
}

func ReadConfig() *Config {
	// TODO: Command line arguments, env config and config file
	var conf Config
	f, e := os.ReadFile("config.yml")
	if e != nil {
		log.Fatalln(e)
		return nil
	}
	yaml.Unmarshal(f, &conf)
	return &conf
	/*v := viper.New()
	v.SetConfigName("blizzard")
	v.SetConfigType("yaml")
	for key, val := range DefaultConfig {
		v.SetDefault(key, val)
	}
	if e := v.ReadInConfig(); e != nil {
		log.Fatalln(e)
	}
	v.SetEnvPrefix("ARCTIC")
	v.AutomaticEnv()
	v.Unmarshal(&conf)
	return*/
}
