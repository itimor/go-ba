package config

import (
	"fmt"

	"github.com/pelletier/go-toml"
)

const (
	TomlConfigFilePath = "./config/test.toml"
)

var (
	Tomlconfig = ParseToml()
)

// toml 不需要映射 struct
func ParseToml() *toml.Tree {
	config, err := toml.LoadFile(TomlConfigFilePath)

	if err != nil {
		fmt.Println("TomlErr ", err.Error())
	}

	return config
}
