package utils

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func InitLogger(filePath string, logger *log.Logger, level log.Level) error {
	DirCreate(filepath.Dir(filePath), 0775)
	FileCreate(filePath)
	var err error
	logf, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		//fmt.Println(err)
		return err
	}

	logger.SetOutput(logf)

	if level < log.PanicLevel && level > log.TraceLevel {
		logger.SetLevel(log.InfoLevel)
	} else {
		logger.SetLevel(level)
	}

	logger.SetReportCaller(true)
	logger.SetFormatter(&log.JSONFormatter{PrettyPrint: false})

	return err
}

func InitLoggerStdout(logger *log.Logger, level log.Level) {
	logger.SetOutput(os.Stdout)

	if level < log.PanicLevel && level > log.TraceLevel {
		logger.SetLevel(log.InfoLevel)
	} else {
		logger.SetLevel(level)
	}

	logger.SetReportCaller(true)
	logger.SetFormatter(&log.JSONFormatter{PrettyPrint: false})
}
