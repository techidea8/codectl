package configure

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type ConsulConfigure struct {
	endpoint string
	key      string
	provider ProviderType
}

func NewConsulConfigure(endpoint string, key string) *ConsulConfigure {
	return &ConsulConfigure{
		endpoint: endpoint,
		key:      key,
		provider: ConsulProvider,
	}
}
func (c *ConsulConfigure) Provider() ProviderType {
	return c.provider
}

// 解析
func (c *ConsulConfigure) Parse(ptr any) (err error) {
	runtimeConfig := viper.New()
	pt := reflect.ValueOf(ptr)
	if pt.Kind() != reflect.Ptr {
		err = errors.New("数据格式不正确")
		return
	}
	key := c.key
	runtimeConfig.AddRemoteProvider(string(c.provider), c.endpoint, key)
	if strings.HasSuffix(key, ".yml") || strings.HasSuffix(key, ".yaml") {
		runtimeConfig.SetConfigType("yaml")
	} else if strings.HasSuffix(key, ".json") {
		runtimeConfig.SetConfigType("json")
	} else if strings.HasSuffix(key, ".xml") {
		runtimeConfig.SetConfigType("xml")
	}

	err = runtimeConfig.ReadRemoteConfig()
	if err != nil {
		return err
	}
	runtimeConfig.Unmarshal(ptr)
	runtimeConfig.WatchRemoteConfigOnChannel()
	runtimeConfig.OnConfigChange(func(in fsnotify.Event) {
		runtimeConfig.Unmarshal(ptr)
		fmt.Println("runtimeConfig.OnConfigChange")
	})
	return
}
