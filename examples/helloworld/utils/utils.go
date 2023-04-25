package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ReadConfig(configFile string, cfg any) error {
	d, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(d, cfg); err != nil {
		return err
	}
	return nil
}
