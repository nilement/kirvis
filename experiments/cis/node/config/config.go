package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/nilement/komrade/experiment"
)

type Config struct {
	Experiments []experiment.Experiment
	ExperimentMap map[string]experiment.Experiment
}

func ReadConfig(filename string) (*Config, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	c.formatMap()
	return c, nil
}

func(c *Config) formatMap() {
	if c.ExperimentMap == nil {
		c.ExperimentMap = make(map[string]experiment.Experiment)
	}
	for _, exp := range c.Experiments {
		c.ExperimentMap[exp.Key] = exp
	}
}