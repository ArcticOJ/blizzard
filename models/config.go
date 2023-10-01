package models

type (
	BlizzardConfig struct {
		Address    `yaml:",inline"`
		PrivateKey string
		Debug      bool `yaml:"-"`
		EnableCORS bool
		RateLimit  uint32 `yaml:"rateLimit"`
		Database   *DatabaseConfig
		Storage    *StorageConfig
		OAuth      map[string]*OAuthProvider
		Redis      *Address
		RabbitMQ   *RabbitMQConfig
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
