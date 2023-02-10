package pkg

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadMonitorDir() []string {
	// @title ReadMonitorDir
	// @param
	// return

	//config settings
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	//
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("read config file failed:", err)
	}
	fmt.Println(viper.GetStringSlice("monitor.dir"))
	return viper.GetStringSlice("monitor.dir")

}
