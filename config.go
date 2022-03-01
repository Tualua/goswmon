package main

import (
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Sites []struct {
		Name            string `yaml:"name"`
		DhcpServerType  string `yaml:"dhcp_server_type"`
		DhcpServer      string `yaml:"dhcp_server"`
		DhcpApiPort     int    `yaml:"dhcp_api_port"`
		Community       string `yaml:"community"`
		DhcpApiLogin    string `yaml:"login"`
		DhcpApiPassword string `yaml:"password"`
	} `yaml:"sites"`
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
