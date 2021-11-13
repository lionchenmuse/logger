package logger

import (
	"fmt"
	"testing"
)

func TestFileLogger(t *testing.T) {
	config := make(map[string]string)
	config["log_path"] = "/home/lionchen/workspace/go-study/logger/log/"
	config["log_name"] = "test3"
	config["log_split_type"] = "size"
	config["log_level"] = "debug"

	logger, err := NewFileLogger(config)
	if err != nil {
		panic(fmt.Sprintf("NewFileLogger error: %v", err))
	}

	logger.Init()

	logger.Debug("user id[%d] come from China", 123)
	logger.Warn("test warn log")
	logger.Fatal("test fatal log")
	logger.Close()
}

func TestConsoleLogger(t *testing.T) {
	config := make(map[string]string)
	config["log_level"] = "error"

	logger, err := NewConsoleLogger(config)
	if err != nil {
		panic(fmt.Sprintf("NewConsoleLogger err: %v", err))
	}
	logger.Debug("user id[%d] come from China", 123)
	logger.Warn("test warn log")
	logger.Fatal("test fatal log")
	logger.Close()

}

func TestLogger(t *testing.T) {
	config := make(map[string]string)
	config["log_path"] = "/home/lionchen/workspace/go-study/logger/log/"
	config["log_name"] = "test4"
	config["log_split_type"] = "size"
	config["log_level"] = "debug"

	InitLogger("file", config)
	Debug("user id[%d] come from China", 123)
	Warn("test warn log")
	Fatal("test fatal log")
	Close()
}
