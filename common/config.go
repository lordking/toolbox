package common

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/viper"
)

// InitConfig 读取配置文件，如果有配置文件(cfgFile != "")，读取配置文件。
// 否则，在/etc/[appName]/、$HOME/.[appName]/、./config/目录下需找文件名为config的配置文件。
func InitConfig(appName, cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(fmt.Sprintf("/etc/%s/", appName))
		viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", appName))
		viper.AddConfigPath("./config/")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Read config file error: ", err)
		os.Exit(0)
	}
}

// ReadConfigFromKey 从读取的配置中寻找符合configKey的配置，并赋值到配置对象(config)定义的变量中。
// 如果config有环境变量定义，那么以环境变量设置为准
func ReadConfigFromKey(configKey string, config interface{}) {

	var err error

	if configKey == "" {
		err = viper.Unmarshal(config)
	} else {
		err = viper.UnmarshalKey(configKey, config)
	}
	CheckFatal(err)

	err = ReadEnv(config)
	CheckFatal(err)
}

//ReadEnv 读取struct内env的标签内容，转为读取环境变量
func ReadEnv(obj interface{}) error {

	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	} else {
		NewError(ErrCodeNotFound, "The object is null")
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		key := field.Tag.Get("env")
		value := os.Getenv(key)

		if key != "" && value != "" {
			val.Field(i).SetString(value)
		}

	}

	return nil
}
