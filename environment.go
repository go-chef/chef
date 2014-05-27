package chef

import (
	"github.com/mitchellh/mapstructure"
)

// Node has a Reader, hey presto
type Environment struct {
	*Reader
	*nativeEnvironment
}

// type RunList []string

type Cookbook map[string]interface{}

// NativeNode represents the native Go version of the deserialized Node type
type nativeEnvironment struct {
	Name            string                 `mapstructure:"name"`
	Environment     string                 `mapstructure:"chef_environment"`
	Default         map[string]interface{} `mapstructure:"default_attributes"`
	Override        map[string]interface{} `mapstructure:"override_attributes"`
	Cookbook        Cookbook               `mapstructure:"cookbook"`
	CookbookVersion Cookbook               `mapstructure:"cookbook_versions"`
}

// NewNode wraps a Node around a pointer to a Reader
func NewEnvironment(reader *Reader) (*Environment, error) {
	environment := Environment{reader, &nativeEnvironment{}}
	if err := mapstructure.Decode(reader, environment.nativeEnvironment); err != nil {
		return nil, err
	}
	return &environment, nil
}
