package main

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Service struct {
		Listen string `yaml:"listen"`
		Port   string `yaml:"port"`
	} `yaml:"service"`
	Netbox struct {
		Address    string `yaml:"address"`
		Token      string `yaml:"token"`
		SwitchRole string `yaml:"switchrole"`
		RackRole   string `yaml:"rackrole"`
	} `yaml:"netbox"`
	Sites []Site `yaml:"sites"`
}

func (s *Config) FindByName(name string) (num int, err error) {
	num = -1
	for i, v := range s.Sites {
		if v.Name == name {
			num = i
		}
	}
	if num == -1 {
		err = fmt.Errorf("site %s not found", name)
	}
	return
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil

}
