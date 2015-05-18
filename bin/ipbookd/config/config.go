package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Host string
	Port int
}

var defaultConfig = Config{"127.0.0.1", 7000}

func GetDefault() Config {
	return defaultConfig
}

func Load(args map[string]interface{}) Config {
	viper.SetConfigName("ipbookd")
	viper.AddConfigPath("/etc/ipbook/")
	viper.AddConfigPath("$HOME/.config/ipbook/")

	viper.ReadInConfig()

	viper.SetEnvPrefix("ipbook")
	viper.AutomaticEnv()

	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("port", "7000")

	setArg("host", args["--host"])
	setArg("port", args["--port"])

	config := Config{}
	config.Host = viper.GetString("host")
	config.Port = viper.GetInt("port")

	return config
}

func setArg(key string, arg interface{}) {
	if nil != arg {
		viper.Set(key, arg.(string))
	}
}
