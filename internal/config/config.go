package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Spec Spec `yaml:"spec"`
}

type Spec struct {
	KubernetesConfig KubernetesConfig `yaml:"kubernetesConfig"`
	FargateConfig    FargateConfig    `yaml:"fargateConfig"`
}

type KubernetesConfig struct {
	Context   *string `yaml:"context,omitempty"`
	Namespace string  `yaml:"namespace,omitempty"`
	Service   string  `yaml:"service,omitempty"`
}

type FargateConfig struct {
	StackName   string `yaml:"name"`
	AccountID   string `yaml:"accountID"`
	Region      string `yaml:"region"`
	ServiceName string `yaml:"serviceName"`
	DomainName  string `yaml:"domainName"`
}

func ParseConfig(filepath string) (*Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
