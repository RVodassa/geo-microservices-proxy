package config

type Config struct {
	GeoServiceAddress  string `env:"GEO_SERVICE_ADDRESS"`
	AuthServiceAddress string `env:"AUTH_SERVICE_ADDRESS"`
	UserServiceAddress string `env:"USER_SERVICE_ADDRESS"`
	HTTPPort           string `env:"HTTP_PORT"`
}

func Load() (*Config, error) {
	return &Config{}, nil
}
