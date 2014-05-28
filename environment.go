package chef

import (
	"github.com/mitchellh/mapstructure"
)

// Environment has a Reader, hey presto
type Environment struct {
	*Reader
	*nativeEnvironment
}

type Cookbook map[string]interface{}

// NativeEnvironment represents the native Go version of the deserialized Environment type
type nativeEnvironment struct {
	Name            string                 `mapstructure:"name"`
	Environment     string                 `mapstructure:"chef_environment"`
	Default         map[string]interface{} `mapstructure:"default_attributes"`
	Override        map[string]interface{} `mapstructure:"override_attributes"`
	Cookbook        Cookbook               `mapstructure:"cookbook"`
	CookbookVersion Cookbook               `mapstructure:"cookbook_versions"`
}

// NewEnvironment wraps a Environment around a pointer to a Reader
func NewEnvironment(reader *Reader) (*Environment, error) {
	environment := Environment{reader, &nativeEnvironment{}}
	if err := mapstructure.Decode(reader, environment.nativeEnvironment); err != nil {
		return nil, err
	}
	return &environment, nil
}
