package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

func Str(name string) string {
	return viper.GetString(name)
}

func Int(name string) int {
	return viper.GetInt(name)
}

func Int64(name string) int64 {
	return viper.GetInt64(name)
}

func Bool(name string) bool {
	return viper.GetBool(name)
}

func Set(name, val string) {
	viper.Set(name, val)
}

func Configure() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(true)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Info().Msg("Unable to read config file")
	}
	for key, value := range viper.AllSettings() {
		fmt.Println(key, value)
	}
}

func TelegramToken() string {
	return viper.GetString("telegram.token")
}

func TelegramAdminId() int64 {
	return viper.GetInt64("telegram.admin")
}

func TelegramAdminLogin() string {
	return viper.GetString("telegram.adminlogin")
}
