package log

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/lordking/toolbox/common"
)

//SetLevel 设置日志级别
type (

	//Defaults 日志配置
	Defaults struct {
		Level      string `json:"level"`
		File       string `json:"file"`
		Production bool   `json:"production"`
	}
)

//SetLogDefaults ...
func SetLogDefaults(configPath string) {

	//读取配置文件
	defaults := new(Defaults)
	err := common.ReadConfig(defaults, configPath)
	common.CheckFatal(err)

	//日志级别
	level, _ := logrus.ParseLevel(defaults.Level)
	fmt.Printf("level:%d\n", level)
	logrus.SetLevel(level)

	//日志文件
	if defaults.Production {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05.000",
		})
	}
}

func Debug(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func Panic(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}
