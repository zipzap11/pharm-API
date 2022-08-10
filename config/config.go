package config

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.AutomaticEnv()
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
	cfg := viper.GetString("postgres.host")
	if cfg == "" {
		return os.Getenv("postgres.host")
	}
	return cfg
}

func GetDBPort() string {
	cfg := viper.GetString("postgres.port")
	if cfg == "" {
		return os.Getenv("postgres.port")
	}
	return cfg
}

func GetDBName() string {
	cfg := viper.GetString("postgres.database")
	if cfg == "" {
		return os.Getenv("postgres.database")
	}
	return cfg
}

func GetDBPassword() string {
	cfg := viper.GetString("postgres.password")
	if cfg == "" {
		return os.Getenv("postgres.password")
	}
	return cfg
}

func GetDBUser() string {
	cfg := viper.GetString("postgres.user")
	if cfg == "" {
		return os.Getenv("postgres.user")
	}
	return cfg
}

func GetRefreshTokenDuration() time.Duration {
	cfg := viper.GetDuration("session.refresh_token_duration")
	if cfg < 0 {
		return time.Minute * 10
	}
	return cfg
}

func GetAccessTokenDuration() time.Duration {
	cfg := viper.GetDuration("session.access_token_duration")
	if cfg < 0 {
		return time.Hour * 96
	}
	return cfg
}

func GetSymmetricKey() string {
	cfg := viper.GetString("session.symmetric_key")
	if cfg == "" {
		return os.Getenv("session.symmetric_key")
	}
	return cfg
}

func GetProvinceAPIUrl() string {
	cfg := viper.GetString("rajaongkir.province")
	if cfg == "" {
		return os.Getenv("rajaongkir.province")
	}
	return cfg
}

func GetStateAPIUrl() string {
	cfg := viper.GetString("rajaongkir.state")
	if cfg == "" {
		return os.Getenv("rajaongkir.state")
	}
	return cfg
}

func GetROAPIKey() string {
	cfg := viper.GetString("rajaongkir.key")
	if cfg == "" {
		return os.Getenv("rajaongkir.key")
	}
	return cfg
}

func GetROPriceURL() string {
	cfg := viper.GetString("rajaongkir.price")
	if cfg == "" {
		return os.Getenv("rajaongkir.price")
	}
	return cfg
}

func GetShippingOrigin() string {
	cfg := viper.GetString("rajaongkir.origin")
	if cfg == "" {
		return os.Getenv("rajaongkir.origin")
	}
	return cfg
}

func GetMidtransAPIKey() string {
	cfg := viper.GetString("midtrans.sandbox_apikey")
	if cfg == "" {
		return os.Getenv("midtrans.sandbox_apikey")
	}
	return cfg
}
