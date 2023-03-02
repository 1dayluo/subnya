// file_logger.go
package logutil

import (
	"fmt"
	"log"
	"os"
)

type FileLogger struct {
	Logger *log.Logger
	File   *os.File
}

func NewFileLogger(filename string) (*FileLogger, error) {
	//@title NewFileLogger:
	//@param
	//Return

	fmt.Println(filename)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	logger := log.New(file, "", log.LstdFlags)
	logger.SetFlags(log.Ldate | log.Lshortfile)
	return &FileLogger{
		Logger: logger,
		File:   file,
	}, nil
}

func (l *FileLogger) Close() error {
	//@title Close:
	//@param
	//Return
	return l.Close()
}
