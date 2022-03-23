package app

import (
	. "go.uber.org/fx"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var Module = Provide(
	RegisterAppConfigs,
)

func RegisterAppConfigs() *viper.Viper {
	_, b, _, _ := runtime.Caller(0)
	path := filepath.Dir(b)

	if err := viper.BindEnv("app.env", "APP_ENV"); err != nil {
		panic(err)
	}

	if !viper.IsSet("app.env") {
		panic("The APP_ENV variable must be set!")
	}

	viper.SetConfigName(viper.GetString("app.env"))
	viper.SetConfigType("yaml")

	viper.AddConfigPath(path)

	err := viper.MergeInConfig()

	if err != nil {
		panic(err)
	}

	return viper.GetViper()
}
