package models

type (
	BlizzardConfig struct {
		Host       string
		PrivateKey string
		Port       uint16
		Debug      bool
		EnableCORS bool
		RateLimit  float64
		Database   *DatabaseConfig
		Storage    *StorageConfig
		Judges     map[string]*Judge
		OAuth      map[string]*OAuthProvider
	}

	OAuthProvider struct {
		ClientId     string
		ClientSecret string
	}

	StorageConfig struct {
		Problems string `yaml:"problems"`
		Posts    string `yaml:"posts"`
		READMEs  string `yaml:"readmes"`
	}

	DatabaseConfig struct {
		Address  string
		Username string
		Password string
		Name     string
		Secure   bool
	}

	Judge struct {
		Address string
		Key     string
	}
)
