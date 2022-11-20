package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Schema struct {
	Name       string            `yaml:"name"`
	Namespaces []NamespaceSchema `yaml:"namespaces"`
}

type NamespaceSchema struct {
	Name      string   `yaml:"name"`
	Relations []string `yaml:"relations"`
}

func ParseSchema() (*Schema, error) {
	s := &Schema{}
	yamlFile, err := ioutil.ReadFile("rbac-schema.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
