package filehandler

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Yaml struct {
	m sync.Mutex
}

func (y *Yaml) Write(data any, file string) error {
	y.m.Lock()
	defer y.m.Unlock()

	yaml, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed at marshal yaml: %v", err)
	}

	err = Utils.WriteFile(file, yaml)
	if err != nil {
		return fmt.Errorf("failed at write to yaml file: %v", err)
	}

	return nil
}

func (y *Yaml) Read(data any, file string) error {
	y.m.Lock()
	defer y.m.Unlock()

	configFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed at open yaml file: %v", err)
	}

	yamlParser := yaml.NewDecoder(configFile)
	if err = yamlParser.Decode(data); err != nil {
		return fmt.Errorf("failed at reading from yaml file: %v", err)
	}

	return nil
}

func (y *Yaml) GetExtension() string {
	return "yaml"
}
