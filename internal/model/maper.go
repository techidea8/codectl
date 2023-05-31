package model

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var DataTypeMap map[string]string = map[string]string{
	"bigint":    "uint",
	"bit":       "bool",
	"char":      "string",
	"date":      "restgo.Date",
	"datetime":  "restgo.DateTime",
	"decimal":   "float64",
	"double":    "float64",
	"float":     "float64",
	"int":       "int",
	"integer":   "int",
	"longtext":  "string",
	"mediumint": "int",
	"numeric":   "float64",
	"smallint":  "int",
	"text":      "string",
	"timestamp": "uint",
	"tinyint":   "int",
	"varchar":   "string",
}

func GuesDataType(input string) string {
	if r, ok := DataTypeMap[input]; ok {
		return r
	} else {
		return input
	}
}

func InitConfig() {
	tmp := map[string]string{}
	config := viper.New()
	config.SetConfigType("yml") // 文件类型
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("找不到配置文件..")
		} else {
			fmt.Println("配置文件出错..")
		}
	}
	config.Unmarshal(&tmp)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		config.Unmarshal(&tmp)
		for k, v := range tmp {
			DataTypeMap[k] = v
			fmt.Println(k, "==>", v)
		}

	})
	for k, v := range tmp {
		DataTypeMap[k] = v
		fmt.Println(k, "==>", v)
	}
}
