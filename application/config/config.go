package config

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type Configuration struct {
	Development bool
	ListenAddr  string
	SecretKey   string

	TemplatePath string
	StaticPath   string

	Database struct {
		Provider   string
		ConnString string
	}
}

func Read(path string) (*Configuration, error) {
	config := &Configuration{}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, config)
	return config, err
}
