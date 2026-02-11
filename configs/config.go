package configs

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server struct {
		ENV     Env    `envconfig:"APP_ENV" required:"true"`
		BASEURL string `envconfig:"APP_BASE_URL" required:"true"`
		NAME    string `envconfig:"APP_NAME" required:"true"`
		PORT    string `envconfig:"APP_PORT" required:"true"`
		DEBUG   bool   `envconfig:"APP_DEBUG" default:"false"`
	}

	Database struct {
		Host               string `envconfig:"DB_HOST" required:"true"`
		Port               string `envconfig:"DB_PORT" required:"true"`
		Username           string `envconfig:"DB_USERNAME" required:"true"`
		Password           string `envconfig:"DB_PASSWORD" required:"true"`
		DBName             string `envconfig:"DB_NAME" required:"true"`
		MinIdleConnections int    `envconfig:"DB_MIN_IDDLE_CONN" required:"true"`
		MaxOpenConnections int    `envconfig:"DB_MAX_OPEN_CONN" required:"true"`
	}

	Redis struct {
		Host     string `envconfig:"REDIS_HOST" required:"true"`
		Port     string `envconfig:"REDIS_PORT" required:"true"`
		Password string `envconfig:"REDIS_PASSWORD" required:"true"`
		DB       int    `envconfig:"REDIS_DB" required:"true"`
	}

	JWT struct {
		SigningKey                       string `envconfig:"JWT_SIGNING_KEY" required:"true"`
		TokenExpiration                  int    `envconfig:"JWT_TOKEN_EXPIRATION" required:"true"`
		RefreshTokenExpirationDay        int    `envconfig:"JWT_REFRESH_TOKEN_EXPIRATION_DAY" required:"true"`
		EncryptionKey                    string `envconfig:"JWT_ENCRYPTION_KEY" required:"true"`
		AccountActivationTokenExpiration string `envconfig:"ACCOUNT_ACTIVATION_TOKEN_EXPIRATION" required:"true"`
		ForgotPasswordTokenExpiration    int    `envconfig:"FORGOT_PASSWORD_TOKEN_EXPIRATION" required:"true"`
	}

	Logger struct {
		Stdout        bool     `envconfig:"LOGGER_STDOUT"`
		FileLocation  string   `envconfig:"LOGGER_FILE_LOCATION"`
		FileMaxAge    int      `envconfig:"LOGGER_FILE_MAX_AGE"`
		Level         int8     `envconfig:"LOGGER_LEVEL"`
		Masking       bool     `envconfig:"LOGGER_MASKING"`
		MaskingParams []string `envconfig:"LOGGER_MASKING_PARAMS"`
	}
}

// LoadTest loads test config
func LoadTest() *Config {
	return load()
}

// LoadDefault loads default config from environment variables
func LoadDefault() *Config {
	return load()
}

// load config from environment variables
func load() *Config {
	var c Config

	_ = godotenv.Load() // Load .env file if it exists

	if err := envconfig.Process("", &c); err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	return &c
}
