/*
 * @Author: 1dayluo
 * @Date: 2023-03-02 20:50:18
 * @LastEditTime: 2023-03-11 23:14:22
 */
// logger.go
package logutil

import (
	"fmt"
	"log"
	"os"

	"github.com/1dayluo/subnya/pkg/readconf"
)

func init() {

	os.MkdirAll(readconf.ReadSettingsConfig("logdir"), os.ModePerm)
}

var (
	DefaultLogger = log.New(os.Stderr, "", log.LstdFlags)
)

func SetLogger(l *log.Logger) {
	DefaultLogger = l
}
func Logf(format string, v ...interface{}) {
	DefaultLogger.Printf(format, v...)
}

func Init() error {
	// savePath := readconf.ReadSettingsConfig("logfile") + time.Now().Format("2006-01-02") + ".log"
	savePath := readconf.ReadSettingsConfig("logdir") + "monitor_run.log"
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		homedir, _ := os.UserHomeDir()
		logPath := fmt.Sprintf("%v/.config/subnya/log", homedir)
		logPath = readconf.SetSettingsConfig("logdir", logPath)
		os.MkdirAll(logPath, os.ModePerm)
	}
	logger, err := NewFileLogger(savePath)
	if err != nil {
		return err
	}
	SetLogger(logger.Logger)
	return nil
}
