package config

import (
	"os"

	"github.com/spf13/viper"
)

var Viper *config

type config struct {
	*viper.Viper
}

func init() {
	v := viper.New()
	_env := os.Getenv("APP_ENV_PROFILE")
	if len(_env) == 0 {
		_env = "dev"
	}
	v.SetDefault("env", _env)
	v.BindEnv("env")
	env := v.GetString("env")
	Viper = &config{initConfig(v, env)}
}

func initConfig(v *viper.Viper, env string) *viper.Viper {
	if v == nil {
		v = viper.New()
		v.SetDefault("env", env)
	}
	configFileName := "config-" + env
	v.SetConfigName(configFileName)
	v.AddConfigPath("conf/")
	v.SetConfigType("yaml")
	v.AddConfigPath("./conf/")
	v.AddConfigPath("../conf/")
	v.AddConfigPath("../../conf/")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	/*
	 Note: Use v.MergeConfig if this service
	 has to read from more than 1 config file
	*/

	return v
}
