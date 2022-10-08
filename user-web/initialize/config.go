package initialize

import (
	"Completed_Api/user-web/global"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Config() {
	Viper := viper.New()
	Viper.SetConfigFile("config.yaml")
	err := Viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	/**
	  使用结构体进行接收
	*/
	if err := Viper.Unmarshal(global.ServConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("Config init success: %s", global.ServConfig)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file channed:", e.Name)
		_ = Viper.ReadInConfig()
		_ = Viper.Unmarshal(global.ServConfig)
		fmt.Println(global.ServConfig)
	})
}
