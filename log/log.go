package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	logLevelName = map[int64]string{
		0: "DEBUG",
		1: "INFO",
		2: "WARNING",
		3: "ERROR",
		4: "CRITICAL",
	}

	logLevelNameIndent = map[int64]string{
		0: "DEBUG   ",
		1: "INFO    ",
		2: "WARNING ",
		3: "ERROR   ",
		4: "CRITICAL",
	}
)

const (
	LOG_LEVEL_DEBUG    int64 = 0
	LOG_LEVEL_INFO     int64 = 1
	LOG_LEVEL_WARNING  int64 = 2
	LOG_LEVEL_ERROR    int64 = 3
	LOG_LEVEL_CRITICAL int64 = 4

	LOG_DEPTH_DEFAULT int = 3
)

var logContext context.Context
var outputLevel int64 = LOG_LEVEL_DEBUG

func SetContext(ctx context.Context) {
	logContext = ctx
}

func SetOutputLevel(i int64) {
	outputLevel = i
}

func Debugf(format string, args ...interface{}) {
	if outputLevel <= LOG_LEVEL_DEBUG {
		Logf(LOG_LEVEL_DEBUG, format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if outputLevel <= LOG_LEVEL_INFO {
		Logf(LOG_LEVEL_INFO, format, args...)
	}
}

func Warningf(format string, args ...interface{}) {
	if outputLevel <= LOG_LEVEL_WARNING {
		Logf(LOG_LEVEL_WARNING, format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if outputLevel <= LOG_LEVEL_ERROR {
		Logf(LOG_LEVEL_ERROR, format, args...)
	}
}

func Criticalf(format string, args ...interface{}) {
	if outputLevel <= LOG_LEVEL_CRITICAL {
		Logf(LOG_LEVEL_CRITICAL, format, args...)
	}
}

func LogfWithDepth(
	level int64,
	depth int,
	format string,
	args ...interface{}) {

	pc, _, _, _ := runtime.Caller(depth)
	f := runtime.FuncForPC(pc)
	fullname := f.Name()
	funcname := strings.Split(fullname, ".")[len(strings.Split(fullname, "."))-1]

	s := fmt.Sprintf(format, args...)
	s = strings.TrimRight(s, "\n")
	s = os.Getenv("PROCESS_ID") + " " + logLevelNameIndent[level] + ":" + funcname + "() " + s

	logPrint(s)
}

func Logf(level int64, format string, args ...interface{}) {
	LogfWithDepth(level, LOG_DEPTH_DEFAULT, format, args...)
}

func logPrint(s string) {
	log.Print(s)
}
