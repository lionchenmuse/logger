package logger

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileLogger struct {
	level         int
	logPath       string
	logName       string
	file          *os.File
	warnFile      *os.File
	LogDataChan   chan *LogData
	endChan       chan struct{}
	logSplitType  int
	logSplitSize  int64
	lastSplitHour int
}

func NewFileLogger(config map[string]string) (log LogInterface, err error) {
	logPath, ok := config["log_path"]
	if !ok {
		err = fmt.Errorf("not found log_path ")
		return
	}
	logName, ok := config["log_name"]
	if !ok {
		err = fmt.Errorf("not found log_name ")
		return
	}
	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not found log_level ")
		return
	}

	logChanSize, ok := config["log_chan_size"]
	if !ok {
		logChanSize = "5000"
	}

	var logSplitType int = Split_By_Hour
	var logSplitSize int64
	splitType, ok := config["log_split_type"]
	if !ok {
		splitType = "hour"
	} else {
		if splitType == "size" {
			logSplitSizeConfig, ok := config["log_split_size"]
			if !ok {
				logSplitSizeConfig = "104857600"
			}
			logSplitSize, err = strconv.ParseInt(logSplitSizeConfig, 10, 64)
			if err != nil {
				logSplitSize = 104857600
			}
			logSplitType = Split_By_Size
		} else {
			logSplitType = Split_By_Hour
		}
	}

	chanSize, err := strconv.Atoi(logChanSize)
	if err != nil {
		chanSize = 5000
	}

	level := getLogLevel(logLevel)
	log = &FileLogger{
		level:         level,
		logPath:       logPath,
		logName:       logName,
		LogDataChan:   make(chan *LogData, chanSize),
		endChan:       make(chan struct{}),
		logSplitSize:  logSplitSize,
		logSplitType:  logSplitType,
		lastSplitHour: time.Now().Hour(),
	}
	return
}

func (f *FileLogger) Init() {
	fileName := fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err: %v", fileName, err))
	}
	f.file = file

	// 写错误日志和fatal日志的文件
	fileName = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err: %v", fileName, err))
	}
	f.warnFile = file

	go f.writeLogBackground()
}

func (f *FileLogger) SetLevel(level int) {
	if level < Level_Debug || level > Level_Fatal {
		level = Level_Debug
	}
	f.level = level
}

func (f *FileLogger) writeLogBackground() {
	// 从通道中获取日志内容
	for logData := range f.LogDataChan {
		var file *os.File = f.file
		if logData.WarnAndFatal {
			file = f.warnFile
		}

		// 在写入日志前，先检查是否需要拆分，即是否需要创建新的文件存放日志
		f.checkSplitFile(logData.WarnAndFatal)

		// 写入日志
		fmt.Fprintf(file, "%s %s (%s:%s:%d) %s\n", logData.TimeStr, logData.LevelStr,
			logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
	}
	f.endChan <- struct{}{}
}

func (f *FileLogger) checkSplitFile(warnFile bool) {
	if f.logSplitType == Split_By_Hour {
		f.splitFileHour(warnFile)
		return
	}
	f.splitFileSize(warnFile)
}

func (f *FileLogger) splitFileHour(warnFile bool) {
	now := time.Now()
	hour := now.Hour()
	if hour == f.lastSplitHour {
		return
	}

	f.lastSplitHour = hour
	var backupFileName string
	var fileName string

	if warnFile {
		backupFileName = fmt.Sprintf("%s%s.log.wf_%04d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)

		fileName = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	} else {
		backupFileName = fmt.Sprintf("%s%s.log_%04d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)

		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}
	file := f.file
	if warnFile {
		file = f.warnFile
	}
	// 将之前的文件关闭，重新命名
	file.Close()
	os.Rename(fileName, backupFileName)

	// 生成并打开新的日志文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}

	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}

func (f *FileLogger) splitFileSize(warnFile bool) {
	file := f.file
	if warnFile {
		file = f.warnFile
	}

	statInfo, err := file.Stat()
	if err != nil {
		return
	}

	fileSize := statInfo.Size()
	if fileSize <= f.logSplitSize {
		return
	}

	var backupFileName string
	var fileName string

	now := time.Now()
	if warnFile {
		backupFileName = fmt.Sprintf("%s/%s.log.wf_%04d%02d%02d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

		fileName = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%04d%02d%02d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}
	// 将原日志文件关闭，并重命名
	file.Close()
	os.Rename(fileName, backupFileName)

	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}

	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}

func (f *FileLogger) Debug(format string, args ...interface{}) {
	if f.level > Level_Debug {
		return
	}

	logData := writeLog(Level_Debug, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	if f.level > Level_Trace {
		return
	}

	logData := writeLog(Level_Trace, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	if f.level > Level_Info {
		return
	}

	logData := writeLog(Level_Info, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	if f.level > Level_Warn {
		return
	}

	logData := writeLog(Level_Warn, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	if f.level > Level_Error {
		return
	}

	logData := writeLog(Level_Error, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	if f.level > Level_Fatal {
		return
	}

	logData := writeLog(Level_Fatal, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Close() {
	close(f.LogDataChan)
	<-f.endChan
	f.file.Close()
	f.warnFile.Close()
	close(f.endChan)
}
