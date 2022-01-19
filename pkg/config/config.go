package config

import (
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	setup sync.Once
)

const DefaultDebug = true

const DefaultChatterbugAPIToken = ""

type Config struct {
	Args []string

	ChatterbugAPIToken string

	Debug bool
}

func doSetup() {
	// NOTE: We do this first because we want the environment variables
	if err := godotenv.Load(); err != nil {
		logger.WithError(err).Warn("failed to load .env file")
	}

	viper.SetEnvPrefix("chatterbug_leaderboard")
	viper.AutomaticEnv()

	viper.SetDefault("debug", DefaultDebug)
	viper.SetDefault("chatterbug_api_token", DefaultChatterbugAPIToken)

	pflag.Bool("debug", DefaultDebug, "execute in debug mode")
	pflag.String("chatterbug-api-token", DefaultChatterbugAPIToken, "API token to authenticate with Chatterbug with")

	viper.BindPFlag("chatterbug_api_token", pflag.Lookup("chatterbug-api-token"))
}

func New() *Config {
	setup.Do(doSetup)

	return &Config{
		Args:               pflag.Args(),
		ChatterbugAPIToken: viper.GetString("chatterbug_api_token"),
		Debug:              viper.GetBool("debug"),
	}
}
