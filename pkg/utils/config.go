package utils

import (
	"CallFrescoBot/pkg/consts"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(consts.EnvFile)
	viper.AddConfigPath(consts.EnvFileDirectory)
	err := viper.ReadInConfig()
	if err != nil {
	}
	viper.AutomaticEnv()
}

func GetEnvVar(name string) string {
	if !viper.IsSet(name) {
		return ""
	}
	value := viper.GetString(name)
	return value
}
