package config

import (
	"github.com/spf13/viper"
)

type Server struct {
	Port      string
	JWTSecret string
	JWTRealm  string
}

type Database struct {
	Name          string
	User          string
	Password      string
	Port          string
	Host          string
	MigrationsDir string
}

type Config struct {
	Server   Server
	Database Database
}

var App Config

func Load() {
	viper.SetConfigName("application")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	App = Config{
		Server: Server{
			Port:      getStringOrPanic("APP_PORT"),
			JWTSecret: getStringOrPanic("JWT_SECRET"),
			JWTRealm:  getStringOrPanic("JWT_REALM"),
		},

		Database: Database{
			Name:          getStringOrPanic("DATABASE_NAME"),
			Host:          getStringOrPanic("DATABASE_HOST"),
			Port:          getStringOrPanic("DATABASE_PORT"),
			User:          getStringOrPanic("DATABASE_USER"),
			Password:      getStringOrPanic("DATABASE_PASSWORD"),
			MigrationsDir: getStringOrPanic("DATABASE_MIGRATIONS_DIR"),
		},
	}
}
