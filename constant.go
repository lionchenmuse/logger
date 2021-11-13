package logger

import (
	"strings"
)

const (
	Level_Debug = iota
	Level_Trace
	Level_Info
	Level_Warn
	Level_Error
	Level_Fatal
)

const (
	Split_By_Hour = iota
	Split_By_Size
)

func getLevelText(level int) string {
	switch level {
	case Level_Debug:
		return "DEBUG"
	case Level_Trace:
		return "TRACE"
	case Level_Info:
		return "INFO"
	case Level_Warn:
		return "WARN"
	case Level_Error:
		return "ERROR"
	case Level_Fatal:
		return "FATAL"
	}
	return "UNKNOWN"
}

func getLogLevel(level string) int {
	if level == "" {
		panic("WRONG LEVEL TYPE")
	}
	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		return Level_Debug
	case "TRACE":
		return Level_Trace
	case "WARN":
		return Level_Warn
	case "ERROR":
		return Level_Error
	case "FATAL":
		return Level_Fatal
	}
	return Level_Debug
}
