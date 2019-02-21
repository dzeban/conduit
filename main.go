package main

import (
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Initialize config
func init() {
	const (
		programName = "conduit"
		defaultPort = 8080
		defaultDSN  = "postgres://postgres:postgres@localhost/conduit?sslmode=disable"
	)

	viper.SetDefault("Port", defaultPort)
	viper.SetDefault("DSN", defaultDSN)

	pflag.IntP("port", "p", defaultPort, "Listen port")
	pflag.StringP("dsn", "d", defaultDSN, "Data source name (database connection string)")

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
	viper.SetEnvPrefix(programName)
	viper.AutomaticEnv()
}

func main() {
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %s\n", err)
	}

	server, err := NewServer(config)
	if err != nil {
		panic(err)
	}

	log.Printf("Start listening on %d\n", config.Port)
	server.Run()
}
