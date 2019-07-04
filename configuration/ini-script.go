package configuration

import (
	"fmt"
	"github.com/go-ini/ini"
)

const settingFile = "setting.ini"

func initDefaultConfiguration() *ini.File {
	cfg := ini.Empty()
	cfg.ReflectFrom(NewConfiguration())
	cfg.SaveTo(settingFile)
	cfg.Reload()
	return cfg
}

func LoadConfiguration() *Configuration {
	cfgFile, err := ini.Load(settingFile)
	if err != nil {
		fmt.Println(err.Error())
		cfgFile = initDefaultConfiguration()
		fmt.Println("generate setting.ini file...")
	}
	config := NewConfiguration()
	cfgFile.MapTo(config)
	return config
}
