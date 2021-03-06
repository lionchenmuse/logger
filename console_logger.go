package logger

import (
	"fmt"
	"os"
)

type ConsoleLogger struct {
	level int
}

func NewConsoleLogger(config map[string]string) (log LogInterface, err error) {
	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not found log_level")
		return
	}

	level := getLogLevel(logLevel)
	log = &ConsoleLogger{
		level: level,
	}
	return
}

func (c *ConsoleLogger) Init() {

}

func (c *ConsoleLogger) SetLevel(level int) {
	if level < Level_Debug || level > Level_Fatal {
		level = Level_Debug
	}
	c.level = level
}

func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	if c.level > Level_Debug {
		return
	}
	logData := writeLog(Level_Debug, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s:%s:%d) %s\n", logData.TimeStr,
		logData.LevelStr, logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Trace(format string, args ...interface{}) {
	if c.level > Level_Trace {
		return
	}
	logData := writeLog(Level_Trace, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s:%s:%d) %s\n", logData.TimeStr,
		logData.LevelStr, logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	if c.level > Level_Info {
		return
	}
	logData := writeLog(Level_Info, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s:%s:%d) %s\n", logData.TimeStr,
		logData.LevelStr, logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	if c.level > Level_Warn {
		return
	}
	logData := writeLog(Level_Warn, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s:%s:%d) %s\n", logData.TimeStr,
		logData.LevelStr, logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	if c.level > Level_Error {
		return
	}
	logData := writeLog(Level_Error, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s:%s:%d) %s\n", logData.TimeStr,
		logData.LevelStr, logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	if c.level > Level_Fatal {
		return
	}
	logData := writeLog(Level_Fatal, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s:%s:%d) %s\n", logData.TimeStr,
		logData.LevelStr, logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Close() {

}
