/*
 * @Author: 1dayluo
 * @Date: 2023-02-10 10:50:17
 * @LastEditTime: 2023-03-11 21:47:08
 */
package readconf

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

var configPath = "./config"

func init() {
	homedir, _ := os.UserHomeDir()
	configPath = fmt.Sprintf("%v/.config/subnya/", homedir)
	os.MkdirAll(configPath, os.ModePerm)
	ioutil.ReadFile("./config/config.yml")
	bytesRead, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(configPath+"config.yml", bytesRead, 0644)
	if err != nil {
		panic(err)
	}
}
func ReadMonitorDir() []string {
	// @title ReadMonitorDir
	// @param
	// return []string  : return dir list in config.yml

	//config settings
	// viper.SetConfigName("config")
	// viper.AddConfigPath("./config")
	viper.AddConfigPath(configPath)
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
	// viper.SetConfigName("config")
	// viper.AddConfigPath("./config")
	viper.AddConfigPath(configPath)
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
	// viper.SetConfigName("config")
	// viper.AddConfigPath("./config")
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()

	//
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("read config file failed:", err)
	}
	// fmt.Println(viper.GetString(fmt.Sprintf("database.%s", key)))
	return viper.GetString(fmt.Sprintf("sqlite.%s", key))

}

func ReadSettingsConfig(key string) string {
	// @title ReadMonitorDir
	// @param
	// return

	//config settings
	// viper.SetConfigName("config")
	// viper.AddConfigPath("./config")
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()

	//
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("read config file failed:", err)
	}

	// fmt.Println(viper.GetString(fmt.Sprintf("database.%s", key)))
	return viper.GetString(fmt.Sprintf("monitor.settings.%s", key))

}
