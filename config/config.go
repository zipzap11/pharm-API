package config

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			logrus.Error("config file not found")
			logrus.Error(err)
		} else {
			// Config file was found but another error was produced
			logrus.Error(err)
		}
	}
	logrus.Info("successfully load config file")
}

func GetDBHost() string {
	return viper.GetString("postgres.host")
}

func GetDBPort() string {
	return viper.GetString("postgres.port")
}

func GetDBName() string {
	return viper.GetString("postgres.database")
}

func GetDBPassword() string {
	return viper.GetString("postgres.password")
}

func GetDBUser() string {
	return viper.GetString("postgres.user")
}

func GetRefreshTokenDuration() time.Duration {
	return viper.GetDuration("session.refresh_token_duration")
}

func GetAccessTokenDuration() time.Duration {
	return viper.GetDuration("session.access_token_duration")
}

func GetSymmetricKey() string {
	return viper.GetString("session.symmetric_key")
}
