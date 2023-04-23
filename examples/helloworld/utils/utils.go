package utils

import (
	"encoding/json"
	"io"
	"os"
)

func ReadConfig(configFile string, cfg any) error {
	f, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, cfg); err != nil {
		return err
	}
	return nil
}
