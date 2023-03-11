/*
 * @Author: 1dayluo
 * @Date: 2023-03-02 20:50:18
 * @LastEditTime: 2023-03-11 21:53:30
 */
// logger.go
package logutil

import (
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
	logger, err := NewFileLogger(savePath)
	if err != nil {
		return err
	}
	SetLogger(logger.Logger)
	return nil
}
