package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"k8s.io/utils/env"

	"github.com/nilement/kubelet/experiment"
)

type Config struct {
	Experiments    []experiment.Experiment
	ExperimentMap  map[string]experiment.Experiment
	NodeName       string
	KubeconfigPath string
}

func ReadConfig(filename string) (*Config, error) {
	nodeName := env.GetString("NODE_NAME", "")
	if nodeName == "" {
		return nil, fmt.Errorf("node name not found")
	}

	kubeconfig := env.GetString("KUBECONFIG", "")
	if kubeconfig == "" {
		return nil, fmt.Errorf("kubeconfig not found")
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}
	// TODO: fix
	c.NodeName = nodeName
	c.KubeconfigPath = kubeconfig

	c.formatMap()
	return c, nil
}

func (c *Config) formatMap() {
	if c.ExperimentMap == nil {
		c.ExperimentMap = make(map[string]experiment.Experiment)
	}
	for _, exp := range c.Experiments {
		c.ExperimentMap[exp.Key] = exp
	}
}
