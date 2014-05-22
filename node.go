package chef

import (
	"github.com/mitchellh/mapstructure"
)

// Node has a Reader, hey presto
type Node struct {
	*Reader
	*nativeNode
}

type RunList []string

// NativeNode represents the native Go version of the deserialized Node type
type nativeNode struct {
	Name        string                 `mapstructure:"name"`
	RunList     RunList                `mapstructure:"run_list"`
	Environment string                 `mapstructure:"chef_environment"`
	Automatic   map[string]interface{} `mapstructure:"automatic"`
	Normal      map[string]interface{} `mapstructure:"normal"`
	Default     map[string]interface{} `mapstructure:"default"`
	Override    map[string]interface{} `mapstructure:"override"`
}

// NewNode wraps a Node around a pointer to a Reader
func NewNode(reader *Reader) (*Node, error) {
	node := Node{reader, &nativeNode{}}
	if err := mapstructure.Decode(reader, node.nativeNode); err != nil {
		return nil, err
	}
	return &node, nil
}
