package config

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Search string
	Messanger Messanger
}

type Messanger struct {
	Token string
	Ids []int64
}

func Read(path string) (Configuration, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return Configuration{}, fmt.Errorf("could not read configuration file: %w", err)
	}


	var config Configuration
	if err := yaml.NewDecoder(bytes.NewReader(file)).Decode(&config); err != nil {
		return Configuration{}, fmt.Errorf("could not parse configuration: %w", err)
	}

	return config, nil
}

