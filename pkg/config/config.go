package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// NewConfig reads config from configFile
// config must be reference to your config structure
func NewConfig(configFile string, config interface{}) error {
	const fn = "config.NewConfig"

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("%s: unable to read config file: %v", fn, err)
	}

	err = json.Unmarshal(file, config)
	if err != nil {
		return fmt.Errorf("%s: unable to unmarshal config: %v", fn, err)
	}

	return nil
}
