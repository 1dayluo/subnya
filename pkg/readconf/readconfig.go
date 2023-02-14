package readconf

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadMonitorDir() []string {
	// @title ReadMonitorDir
	// @param
	// return []string  : return dir list in config.yml

	//config settings
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	//
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("read config file failed:", err)
	}
	// fmt.Println(viper.GetStringSlice("monitor.dir"))
	return viper.GetStringSlice("monitor.dir")

}

func ReadRedisConfig(key string) string {
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
	// fmt.Println(viper.GetString(fmt.Sprintf("database.%s", key)))
	return viper.GetString(fmt.Sprintf("redis.%s", key))

}

func ReadSqlConfig(key string) string {
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
	// fmt.Println(viper.GetString(fmt.Sprintf("database.%s", key)))
	return viper.GetString(fmt.Sprintf("sqlite.%s", key))

}
