package core

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Host       string            `yaml:"host"`
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

var DefaultConfig = map[string]interface{}{
	"Host":              "0.0.0.0",
	"Port":              2999,
	"Database.Address":  "localhost:5432",
	"Database.Username": "alphanecron",
	"Database.Password": "necronthedev",
	"Database.Name":     "arctic",
}

func ReadConfig() *Config {
	var conf []Config
	f, e := os.ReadFile("config.yml")
	if e != nil {
		log.Fatalln(e)
		return nil
	}
	yaml.Unmarshal(f, &conf)
	return &conf[0]
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
