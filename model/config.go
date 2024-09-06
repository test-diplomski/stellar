package model

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ConfigFile(n ...string) (*Config, error) {
	path := "config.yml"
	if len(n) > 0 {
		path = n[0]
	}

	yamlFile, err := ioutil.ReadFile(path)
	check(err)

	var conf Celestial
	err = yaml.Unmarshal(yamlFile, &conf)
	check(err)

	return &conf.Conf, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Celestial struct {
	Conf Config `yaml:"stellar"`
}

type Config struct {
	ConfVersion    string            `yaml:"version"`
	Address        string            `yaml:"address"`
	Endpoints      []string          `yaml:"db"`
	Syncer         string            `yaml:"syncer"`
	STopic         string            `yaml:"stopic"`
	InstrumentConf map[string]string `yaml:"instrument"`
}
