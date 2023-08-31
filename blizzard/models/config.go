package models

type (
	BlizzardConfig struct {
		Address    `yaml:",inline"`
		PrivateKey string
		Debug      bool
		EnableCORS bool
		RateLimit  float64
		Database   *DatabaseConfig
		Storage    *StorageConfig
		Judges     map[string]string
		OAuth      map[string]*OAuthProvider
		Redis      *RedisConfig
		RabbitMQ   *RabbitMQConfig
	}

	RabbitMQConfig struct {
		Username string
		Password string
		Address  `yaml:",inline"`
	}

	Address struct {
		Host string
		Port uint16
	}

	RedisConfig struct {
		Address `yaml:",inline"`
		DB      int
	}

	OAuthProvider struct {
		ClientID     string
		ClientSecret string
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
