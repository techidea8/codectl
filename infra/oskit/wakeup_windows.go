//go:build windows
// +build windows

package oskit

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

/*
Windows Registry Editor Version 5.00

[HKEY_CLASSES_ROOT\bssoft]
@="bssoft Protocol"
"URL Protocol"=""

[HKEY_CLASSES_ROOT\bssoft\DefaultIcon]
@="C:\\Program Files (x86)\\Internet Explorer\\iexplore.exe"

[HKEY_CLASSES_ROOT\bssoft\shell]
@=""

[HKEY_CLASSES_ROOT\bssoft\shell\open]
@=""

[HKEY_CLASSES_ROOT\bssoft\shell\open\command]
@="C:\\Program Files (x86)\\Internet Explorer\\iexplore.exe http://localhost/plugin/camera/preview?id=%1 "
*/
//schema string 自定义协议

// exepath string 应用程序路径
func Wakeuper(schema string, exepath string, params ...string) (err error) {
	subkey := schema
	key, exists, err := registry.CreateKey(registry.CLASSES_ROOT, subkey, registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	defer key.Close()
	key.SetStringValue("URL Protocol", "")
	subkey = "shell"
	key, _, err = registry.CreateKey(key, subkey, registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	subkey = "open"
	key, _, err = registry.CreateKey(key, subkey, registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	subkey = "command"
	key, _, err = registry.CreateKey(key, subkey, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("创建%s失败", subkey)
	}
	err = key.SetStringValue("", "\""+exepath+"\"")
	return err
}
