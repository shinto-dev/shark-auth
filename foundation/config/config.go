package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	// Enable viper to read env variables
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithError(err).Error("config read failed")
	}
}

func GetString(key string) string {
	fmt.Println(viper.ConfigFileUsed())
	fmt.Println(viper.GetString(key))
	return viper.GetString(key)
}

func GetStringOrDefault(key string, defaultValue string) string {
	viper.SetDefault(key, defaultValue)
	return viper.GetString(key)
}
