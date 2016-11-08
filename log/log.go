package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/lordking/toolbox/common"
	"github.com/spf13/viper"
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
func SetLogDefaults(key string) {

	var err error

	//读取配置文件
	defaults := new(Defaults)
	if key == "" {
		err = viper.Unmarshal(defaults)
	} else {
		err = viper.UnmarshalKey(key, defaults)
	}
	common.CheckFatal(err)

	//日志级别
	level, _ := logrus.ParseLevel(defaults.Level)
	logrus.SetLevel(level)

	//输出日志文件
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

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}
