package main

type RBACSchema struct {
	name      string      `yaml:"name"`
	namespace []Namespace `yaml:"namespaces"`
}

type Namespace struct {
}

func setupRBAC() {

}
