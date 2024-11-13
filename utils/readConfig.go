package utils

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Mappers of config
type Config struct {
	Profile map[string]Profile `yaml:"Profile"`
}

type Profile struct {
	Desc    string            `yaml:"desc"`
	Servers map[string]Server `yaml:"servers"`
}

type Server struct {
	Desc string `yaml:"desc"`
}

func ReadConfig() (Config, error) {
	var config Config
	file, err := os.Open("configuration.yaml")
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
		return config, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	if err := parseYAML(file, &config); err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
		return config, err
	}
	return config, nil
}

func parseYAML(r io.Reader, out interface{}) error {
	decoder := yaml.NewDecoder(r)
	return decoder.Decode(out)
}
