package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Host    string
	Port    int
	Timeout time.Duration
}

var defaultConfig = Config{"127.0.0.1", 7000, 0}

func GetDefault() Config {
	return defaultConfig
}

func Load(args map[string]interface{}) Config {
	viper.SetConfigName("ipbook")
	viper.AddConfigPath("/etc/ipbook/")
	viper.AddConfigPath("$HOME/.config/ipbook/")

	viper.ReadInConfig()

	viper.SetEnvPrefix("ipbook")
	viper.AutomaticEnv()

	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("port", "7000")
	viper.SetDefault("timeout", "0")

	setArg("host", args["--host"])
	setArg("port", args["--port"])
	setArg("timeout", args["--timeout"])

	config := Config{}
	config.Host = viper.GetString("host")
	config.Port = viper.GetInt("port")
	config.Timeout = time.Second * time.Duration(viper.GetInt("timeout"))

	return config
}

func setArg(key string, arg interface{}) {
	if nil != arg {
		viper.Set(key, arg.(string))
	}
}
