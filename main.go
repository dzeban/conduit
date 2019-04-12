package main

import (
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Initialize config
func init() {
	const (
		programName         = "conduit"
		defaultServerPort   = 8080
		defaultArticlesDSN  = "postgres://postgres:postgres@localhost/conduit?sslmode=disable"
		defaultArticlesType = "postgres"
		defaultUsersDSN     = "postgres://postgres:postgres@localhost/conduit?sslmode=disable"
		defaultUsersType    = "postgres"
	)

	viper.SetDefault("Server.Port", defaultServerPort)
	viper.SetDefault("Articles.DSN", defaultArticlesDSN)
	viper.SetDefault("Articles.Type", defaultArticlesType)
	viper.SetDefault("Users.DSN", defaultUsersDSN)
	viper.SetDefault("Users.Type", defaultUsersType)

	pflag.IntP("server.port", "p", defaultServerPort, "Listen port")
	pflag.StringP("articles.dsn", "a", defaultArticlesDSN, "Data source name (database connection string)")
	pflag.StringP("articles.type", "t", defaultArticlesType, "Articles service type")
	pflag.StringP("users.dsn", "d", defaultUsersDSN, "Data source name (database connection string)")
	pflag.StringP("users.type", "u", defaultUsersType, "Users service type")

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

	log.Printf("Start listening on %d\n", config.Server.Port)
	server.Run()
}
