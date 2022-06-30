package config

import (
	"log"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	instance Config
	once     sync.Once
)

type ConfigVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type ConfigFile struct {
	Environments []ConfigVar `yaml:"envs"`
	Secrets      []ConfigVar `yaml:"secrets"`
}

type ConfigKey string

type Config struct {
	values map[ConfigKey]string
}

var (
	PostgresConnection ConfigKey = "pg-dsn"
	Port               ConfigKey = "port"
	MultipartFileKey   ConfigKey = "multipartfilekey"
	RootFileDirectory  ConfigKey = "rootfiledir"
	Domain             ConfigKey = "domain"
)

func (c *Config) setData(data []ConfigVar) {
	if len(c.values) == 0 {
		c.values = make(map[ConfigKey]string, 0)
	}

	for _, val := range data {
		c.values[ConfigKey(strings.ToLower(val.Name))] = val.Value
		log.Printf("%s = %s", val.Name, val.Value)
	}
}

func GetValue(key ConfigKey) string {
	if val := os.Getenv(strings.ToUpper(string(key))); val != "" {
		log.Println("env: ", val)
		return val
	}

	val, ok := instance.values[key]
	if !ok {
		return ""
	}
	return val
}

func uploadConfigFromYaml() ConfigFile {
	valuesFile := os.Getenv("VALUES")
	if valuesFile == "" {
		valuesFile = "./.o3/k8s/values_local.yaml"
	}

	data, err := os.ReadFile(valuesFile)
	if err != nil {
		log.Fatalln(err.Error())
	}
	var config ConfigFile
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return config
}

func newConfig() Config {
	confFile := uploadConfigFromYaml()

	config := Config{}
	config.setData(confFile.Secrets)
	config.setData(confFile.Environments)

	return config
}

func init() {
	once.Do(func() {
		instance = newConfig()
	})
}
