package logger

import (
	"log"
	"os"
)

type Loggers struct {
	LogInfo  *log.Logger
	LogDebug *log.Logger
}

var Log *Loggers

func Run() {
	Log = New()
}

var fileDebugLogName = "logs/log_debug.log"

func New() *Loggers {
	fileDebug, err := os.OpenFile(fileDebugLogName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("--- не удалось создать журнал debug-log ---")
	}
	flags := log.LstdFlags | log.Lshortfile
	logInfo := log.New(os.Stdout, "INFO:\t", log.LstdFlags)
	logDebug := log.New(fileDebug, "Debug:\t", flags)
	return &Loggers{
		LogInfo:  logInfo,
		LogDebug: logDebug,
	}
}

func (l *Loggers) Info(args ...interface{}) {
	l.LogInfo.Println(args...)
}

func (l *Loggers) Debug(args ...interface{}) {
	l.LogDebug.Println(args...)
}
