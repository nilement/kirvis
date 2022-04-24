package config

import (
	"fmt"
	"io/ioutil"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	"github.com/nilement/apiserver/experiment"
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

func ReadAPIServer(filename string) (*v1.Pod, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	pod := &v1.Pod{}
	err = yaml.Unmarshal(buf, pod)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return pod, nil
}

func WriteAPIServer(out string, apiServer *v1.Pod) error {
	bytes, err := yaml.Marshal(apiServer)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(out, bytes, 0644)
}

func(c *Config) formatMap() {
	if c.ExperimentMap == nil {
		c.ExperimentMap = make(map[string]experiment.Experiment)
	}
	for _, exp := range c.Experiments {
		c.ExperimentMap[exp.Key] = exp
	}
}