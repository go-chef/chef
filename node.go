package chef

import "fmt"

type NodeService struct {
	client *Client
}

// Node represents the native Go version of the deserialized Node type
type Node struct {
	Name                string                 `json:"name"`
	Environment         string                 `json:"chef_environment"`
	ChefType            string                 `json:"chef_type"`
	AutomaticAttributes map[string]interface{} `json:"automatic"`
	NormalAttributes    map[string]interface{} `json:"normal"`
	DefaultAttributes   map[string]interface{} `json:"default"`
	OverrideAttributes  map[string]interface{} `json:"override"`
	JsonClass           string                 `json:"json_class"`
	RunList             []string               `json:"run_list"`
}

// List lists the nodes in the Chef server.
//
// Chef API docs: http://docs.opscode.com/api_chef_server.html#id25
func (e *NodeService) List() (data map[string]string, err error) {
	err = e.client.magicRequestDecoder("GET", "nodes", nil, &data)
	return
}

// Get gets a node from the Chef server.
//
// Chef API docs: http://docs.opscode.com/api_chef_server.html#id28
func (e *NodeService) Get(name string) (node Node, err error) {
	url := fmt.Sprintf("nodes/%s", name)
	err = e.client.magicRequestDecoder("GET", url, nil, &node)
	return
}
