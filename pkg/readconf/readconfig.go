/*
 * @Author: 1dayluo
 * @Date: 2023-02-10 10:50:17
 * @LastEditTime: 2023-03-11 23:03:31
 */
package readconf

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/viper"
)
var homedir, _  = os.UserHomeDir()
var ConfigPath = fmt.Sprintf("%v/.config/subnya/", homedir)
var downloadSamplePath = "https://raw.githubusercontent.com/1dayluo/subnya/master/config/config.yml"

func init() {
	/**
	 * @description: Update the subdomain under the domain and check its response code
	 * @return {*}
	 */


	configFile := "config.yml"
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		// get response from github
		response, err := http.Get(downloadSamplePath)
		if err != nil {
			fmt.Println("[*]error,can not download sample config path!")
		}
		defer response.Body.Close()
		// Create dir&file
		os.MkdirAll(ConfigPath, os.ModePerm)
		file, err := os.Create(ConfigPath + configFile)
		if err != nil {
			fmt.Println("[ERR]read config file failed!", err)
		}
		defer file.Close()
		// Copy response text to file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Println("[ERR]read config file failed!", err)
		}

	}

}
func ReadMonitorDir() []string {
	// @title ReadMonitorDir
	// @param
	// return []string  : return dir list in config.yml

	//config settings
	viper.SetConfigName("config")
	// viper.AddConfigPath("./config")
	viper.AddConfigPath(ConfigPath)
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
	// viper.AddConfigPath("./config")
	viper.AddConfigPath(ConfigPath)
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
	// viper.AddConfigPath("./config")
	viper.AddConfigPath(ConfigPath)
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
	viper.SetConfigName("config")
	// viper.AddConfigPath("./config")
	viper.AddConfigPath(ConfigPath)
	viper.AutomaticEnv()

	//
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("read config file failed:", err)
	}

	// fmt.Println(viper.GetString(fmt.Sprintf("database.%s", key)))
	return viper.GetString(fmt.Sprintf("monitor.settings.%s", key))

}

func SetSettingsConfig(key string, value string) string {
	// @title SetSettingsConfig
	// @param
	// return
	viper.SetConfigName("config")
	viper.Set("monitor.settings.%s", key)
	err := viper.WriteConfig()
	if err != nil {
		panic(err)
	}
	return viper.GetString(fmt.Sprintf("monitor.settings.%s", key))
}

func SetSqliteConfig(key, value string) string {
	// @title SetSettingsConfig
	// @param
	// return
	viper.SetConfigName("config")
	configKey := fmt.Sprintf("sqlite.%s", key)
	viper.Set(configKey, value)
	err := viper.WriteConfig()
	if err != nil {
		panic(err)
	}
	return viper.GetString(fmt.Sprintf("sqlite.%s", key))
}
