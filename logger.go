package logger

import "fmt"

var log LogInterface

// InitLogger 形参name的可取值包括：file，console
func InitLogger(name string, config map[string]string) (err error) {
	switch name {
	case "file":
		log, err = NewFileLogger(config)
		log.Init()
	case "console":
		log, err = NewConsoleLogger(config)
		log.Init()
	default:
		err = fmt.Errorf("unsupported logger name: %s", name)
	}
	return
}

func Debug(format string, args ...interface{}) {
	log.Debug(format, args...)
}

func Trace(format string, args ...interface{}) {
	log.Trace(format, args...)
}

func Info(format string, args ...interface{}) {
	log.Info(format, args...)
}
func Warn(format string, args ...interface{}) {
	log.Warn(format, args...)
}
func Error(format string, args ...interface{}) {
	log.Error(format, args...)
}
func Fatal(format string, args ...interface{}) {
	log.Fatal(format, args...)
}

func Close() {
	log.Close()
}
