package main

import (
	"time"

	"github.com/lionchenmuse/logger/logger"
)

func main() {
	initLogger("file", "logs/", "user_server", "debug")
	run()
}

func initLogger(name string, logPath string, logName string, level string) (err error) {
	config := make(map[string]string, 8)
	config["log_path"] = logPath
	config["log_name"] = logName
	config["log_level"] = level
	config["log_split_type"] = "size"
	err = logger.InitLogger(name, config)
	if err != nil {
		return
	}

	logger.Debug("init logger success")
	return
}

func run() {
	for {
		logger.Debug("user server is running, %v", time.Now().Format("2006-01-02 15:04:05.999"))
		time.Sleep(time.Second)
	}
}
