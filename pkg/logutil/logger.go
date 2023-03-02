// logger.go
package logutil

import (
	"log"
	"os"
)

var (
	DefaultLogger = log.New(os.Stderr, "", log.LstdFlags)
)

func SetLogger(l *log.Logger) {
	DefaultLogger = l
}
func Logf(format string, v ...interface{}) {
	DefaultLogger.Printf(format, v...)
}
