package app

import (
	"io/ioutil"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"encoding/json"
	"errors"
)

var ErrNoConfFile       = errors.New("ERR: config file not found")
var ErrConfCorrupted    = errors.New("ERR: invalid config data")

func LoadConfig(configPaths []string) *Config{

	slog.DebugLn("[shared.config.app.lib] ", "Handle log paths: ", configPaths)

	configRaw := func(configPaths []string) (configRaw []byte){

		for _,configPath := range configPaths{

			configRaw, err := ioutil.ReadFile(configPath)

			if err == nil && configRaw != nil {
				slog.DebugLn("[shared.config.app.lib] ", "Loaded config: ", configPath)
				return configRaw
			}
		}

		return nil

	}(configPaths)


	if configRaw == nil {
		panic(ErrNoConfFile)
	}

	config := &Config{}
	if err := json.Unmarshal(configRaw, config); err != nil{
		panic(ErrConfCorrupted)
	}

	return config
}