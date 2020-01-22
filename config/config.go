package config

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// Global configuration variables
var Data Config

// Configuration file directory
var path = "./config/config.yaml"

// Configuration structure
type Config struct {
	Server cServer `yaml:"server"`
	Path cPath `yaml:"path"`
	Mysql cMysql `yaml:"mysql"`
}

// Path for Config
type cPath struct {
	Theme string `yaml:"theme"`
	Work string `yaml:"work"`
}

// Server for Config
type cServer struct {
	Title string `yaml:"title"`
	Addr string `yaml:"addr"`
	Password string `yaml:"password"`
	Protocol string `yaml:"protocol"`
}

type cMysql struct {
	User string `yaml:"user"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
	Charset string `yaml:"charset"`
}

func ParseYaml() {
	if err := Data.Get(); err != nil {
		panic(err)
	}
}

// Get configuration information from a configuration file
func (c *Config) Get() error {
	if f, err := os.Open(path); err != nil {
		return err
	} else if err = yaml.NewDecoder(f).Decode(c); err != nil {
		return err
	}
	return nil
}

func (c *Config) String() string {
	byt,err := json.Marshal(c)
	if err != nil {
		log.Println(err)
	}
	return string(byt)
}

// Returns a new Config with configuration information
func New() (*Config, error) {
	conf := &Config{}
	if err := conf.Get(); err != nil {
		return nil, err
	}
	return conf, nil
}
