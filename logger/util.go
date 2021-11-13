package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

type LogData struct {
	Message      string
	TimeStr      string
	LevelStr     string
	FileName     string
	FuncName     string
	LineNo       int
	WarnAndFatal bool
}

/*
1. 当业务调用打日志的方法时，我们把日志相关的数据写入到chan（队列）
2. 然后我们有一个后台的线程不断的从chan里面获取这些日志，最终写入到文件。
*/
func writeLog(level int, format string, args ...interface{}) *LogData {
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05.999")
	levelStr := getLevelText(level)

	fileName, funcName, lineNo := GetLineInfo()
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	msg := fmt.Sprintf(format, args...)

	logData := &LogData{
		Message:      msg,
		TimeStr:      nowStr,
		LevelStr:     levelStr,
		FileName:     fileName,
		FuncName:     funcName,
		LineNo:       lineNo,
		WarnAndFatal: false,
	}

	if level == Level_Error || level == Level_Warn || level == Level_Fatal {
		logData.WarnAndFatal = true
	}
	return logData
}

func GetLineInfo() (fileName string, funcName string, lineNo int) {
	// Caller：参数从0开始，可以递增
	// 参数0：返回当前函数（即GetLineInfo）的函数指针，该函数所在的文件名，和对应行号
	// 参数1：返回上一级函数（即writeLog）的函数指针，该函数所在的文件名，和对应行号
	// 参数2：返回上上级函数（即logger.(*FileLogger).Debug）的函数指针，所在的文件名（file_logger.go），和对应行号
	// 参数3：返回上上上即函数（即logger.TestFileLogger）的函数指针，所在的文件名（logger_test），和对应行号
	// 依此类推
	pc, file, line, ok := runtime.Caller(4)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNo = line
	}
	return
}
