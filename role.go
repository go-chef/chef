package chef

import (
	"github.com/mitchellh/mapstructure"
)

// Role has a Reader, hey presto
type Role struct {
	*Reader
	*nativeRole
}

// NativeRole represents the native Go version of the deserialized Role type
type nativeRole struct {
	Name     string                 `mapstructure:"name"`
	RunList  RunList                `mapstructure:"run_list"`
	Default  map[string]interface{} `mapstructure:"default_attributes"`
	Override map[string]interface{} `mapstructure:"override_attributes"`
}

// NewRole wraps a Role around a pointer to a Reader
func NewRole(reader *Reader) (*Role, error) {
	role := Role{reader, &nativeRole{}}
	if err := mapstructure.Decode(reader, role.nativeRole); err != nil {
		return nil, err
	}
	return &role, nil
}
