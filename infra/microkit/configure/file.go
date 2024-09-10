package configure

import (
	"errors"
	"reflect"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type FileConfigure struct {
	path     string
	provider ProviderType
}

func NewFileConfigure(path string) *FileConfigure {
	return &FileConfigure{
		path:     path,
		provider: FileProvider,
	}
}
func (c *FileConfigure) Provider() ProviderType {
	return c.provider
}

// 解析
func (c *FileConfigure) Parse(ptr any, key ...string) (err error) {
	runtimeConfig := viper.New()
	pt := reflect.ValueOf(ptr)
	if pt.Kind() != reflect.Ptr {
		err = errors.New("数据格式不正确")
		return
	}
	runtimeConfig.SetConfigFile(c.path)
	path := c.path
	if strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml") {
		runtimeConfig.SetConfigType("yaml")
	} else if strings.HasSuffix(path, ".json") {
		runtimeConfig.SetConfigType("json")
	} else if strings.HasSuffix(path, ".xml") {
		runtimeConfig.SetConfigType("xml")
	}
	err = runtimeConfig.ReadInConfig()
	if err != nil {
		return err
	}
	runtimeConfig.Unmarshal(ptr)
	runtimeConfig.OnConfigChange(func(in fsnotify.Event) {
		runtimeConfig.Unmarshal(ptr)
	})
	return
}
