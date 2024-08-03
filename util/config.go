package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config contains all environment variables.
type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationUrl         string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMYTRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	FromEmailName        string        `mapstructure:"FROM_EMAIL_NAME"`
	FromEmailAddress     string        `mapstructure:"FROM_EMAIL_ADDRESS"`
	FromEmailPassword    string        `mapstructure:"FROM_EMAIL_PASSWORD"`
}

// GetConfig get configs from app.env in path or from environemnt variables.
func GetConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	viper.Unmarshal(&config)

	return
}
