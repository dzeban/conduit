package main

import (
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	PROGRAM_NAME = "conduit"
	DEFAULT_PORT = 8080
	DEFAULT_DSN  = "postgres://postgres:postgres@localhost/conduit?sslmode=disable"
)

// Initialize config
func init() {
	viper.SetDefault("Port", DEFAULT_PORT)
	viper.SetDefault("DSN", DEFAULT_DSN)

	pflag.IntP("port", "p", DEFAULT_PORT, "Listen port")
	pflag.StringP("dsn", "d", DEFAULT_DSN, "Data source name (database connection string)")

	config := pflag.StringP("config", "c", "", "Path to config file")
	pflag.Parse()

	if *config != "" {
		viper.SetConfigFile(*config)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("failed to read config: %s", err)
		}
	}

	viper.BindPFlags(pflag.CommandLine)
	viper.SetEnvPrefix(PROGRAM_NAME)
	viper.AutomaticEnv()
}

func main() {
	var config Config
	viper.Unmarshal(&config)

	server := NewServer(config)
	server.Run()
}
