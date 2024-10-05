package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucax88x/wentsketchy/internal/homedir"
	"gopkg.in/yaml.v2"
)

type Cfg struct {
	Left       []string `yaml:"left"`
	Center     []string `yaml:"center"`
	Right      []string `yaml:"right"`
	LeftNotch  []string `yaml:"left_notch"`
	RightNotch []string `yaml:"right_notch"`
}

func ReadYaml() (*Cfg, error) {
	var cfg Cfg

	dir, err := homedir.Get()

	if err != nil {
		//nolint:errorlint // no wrap
		return nil, fmt.Errorf("config: error getting home dir. %v", err)
	}

	yamlData, err := os.ReadFile(filepath.Join(dir, "config.yaml"))

	if err != nil {
		//nolint:errorlint // no wrap
		return nil, fmt.Errorf("config: could not read file. %v", err)
	}

	err = yaml.Unmarshal(yamlData, &cfg)

	if err != nil {
		//nolint:errorlint // no wrap
		return nil, fmt.Errorf("config: could not unmarshal cfg. %v", err)
	}

	return &cfg, nil
}
