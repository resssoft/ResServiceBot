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
		log.Error().Err(err).Msg("Unable to read config file")
	}
	for key, value := range viper.AllSettings() {
		fmt.Println(key, value)
	}
}

func TelegramToken(botName string) string {
	return viper.GetString("telegram.bots." + botName + ".token")
}

func TelegramAdminId(botName string) int64 {
	return viper.GetInt64("telegram.bots." + botName + ".admin")
}

func TelegramAdminLogin(botName string) string {
	return viper.GetString("telegram.bots." + botName + ".adminlogin")
}

func TelegramIsWebMode(botName string) bool {
	return viper.GetBool("telegram.bots." + botName + ".web")
}

func TelegramBotUrl(botName string) string {
	return viper.GetString("telegram.bots." + botName + ".uri")
}

func TelegramBotCommand(botName string) string {
	return viper.GetString("telegram.bots." + botName + ".command")
}

func WebServerAddr() string {
	return viper.GetString("server.url")
}
