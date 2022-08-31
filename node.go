package chef

import (
	"errors"
	"fmt"
)

var (
	ErrPathNotFound   = errors.New("attributte path not found")
	ErrNoPathProvided = errors.New("no path was provided")
)

type NodeService struct {
	client *Client
}

// Node represents the native Go version of the deserialized Node type
type Node struct {
	Name                string                 `json:"name"`
	Environment         string                 `json:"chef_environment,omitempty"`
	ChefType            string                 `json:"chef_type,omitempty"`
	AutomaticAttributes map[string]interface{} `json:"automatic,omitempty"`
	NormalAttributes    map[string]interface{} `json:"normal,omitempty"`
	DefaultAttributes   map[string]interface{} `json:"default,omitempty"`
	OverrideAttributes  map[string]interface{} `json:"override,omitempty"`
	JsonClass           string                 `json:"json_class,omitempty"`
	//TODO: use the RunList struct for this
	RunList     []string `json:"run_list,omitempty"`
	PolicyName  string   `json:"policy_name,omitempty"`
	PolicyGroup string   `json:"policy_group,omitempty"`
}

// GetAttribute will fetch an attribute from that provided path considering the right attribute precedence.
func (e *Node) GetAttribute(paths ...string) (interface{}, error) {
	if len(paths) <= 0 {
		return nil, ErrNoPathProvided
	}

	// this follows the Chef attribute precedence: https://docs.chef.io/attribute_precedence/
	attrList := []map[string]interface{}{e.AutomaticAttributes, e.OverrideAttributes, e.NormalAttributes, e.DefaultAttributes}

	for _, attrs := range attrList {
		attr, err := lookupAttribute(attrs, paths...)
		if err != nil {
			if errors.Is(err, ErrPathNotFound) {
				continue
			}

			return nil, err
		}

		return attr, nil
	}

	return nil, ErrPathNotFound
}

// looks up a complete path in the provided attribute map.
func lookupAttribute(attrs map[string]interface{}, paths ...string) (interface{}, error) {
	if len(paths) <= 0 {
		return nil, ErrPathNotFound
	}

	currentPath, remainingPaths := paths[0], paths[1:]

	if attr, ok := attrs[currentPath]; ok {
		if len(remainingPaths) <= 0 {
			return attr, nil // we are at the last provided part of the path
		}

		// otherwise keep looking until we reach the end
		return lookupAttribute(attr.(map[string]interface{}), remainingPaths...)
	}

	return nil, ErrPathNotFound
}

type NodeResult struct {
	Uri string `json:"uri"`
}

// NewNode is the Node constructor method
func NewNode(name string) (node Node) {
	node = Node{
		Name:        name,
		Environment: "_default",
		ChefType:    "node",
		JsonClass:   "Chef::Node",
	}
	return
}

// List lists the nodes in the Chef server.
//
// Chef API docs: https://docs.chef.io/api_chef_server.html#nodes
func (e *NodeService) List() (data map[string]string, err error) {
	err = e.client.magicRequestDecoder("GET", "nodes", nil, &data)
	return
}

// Get gets a node from the Chef server.
//
// Chef API docs: https://docs.chef.io/api_chef_server.html#nodes-name
func (e *NodeService) Get(name string) (node Node, err error) {
	url := fmt.Sprintf("nodes/%s", name)
	err = e.client.magicRequestDecoder("GET", url, nil, &node)
	return
}

// Head gets a node from the Chef server. Does not return a json body.
//
// Chef API docs: https://docs.chef.io/api_chef_server.html#nodes-name
func (e *NodeService) Head(name string) (err error) {
	url := fmt.Sprintf("nodes/%s", name)
	err = e.client.magicRequestDecoder("HEAD", url, nil, nil)
	return
}

// Post creates a Node on the chef server
//
// Chef API docs: https://docs.chef.io/api_chef_server.html#nodes
func (e *NodeService) Post(node Node) (data *NodeResult, err error) {
	body, err := JSONReader(node)
	if err != nil {
		return
	}

	err = e.client.magicRequestDecoder("POST", "nodes", body, &data)
	return
}

// Put updates a node on the Chef server.
//
// Chef API docs: https://docs.chef.io/api_chef_server.html#nodes-name
// TODO: We might want to change the name. name and data should be separate structures
func (e *NodeService) Put(n Node) (node Node, err error) {
	url := fmt.Sprintf("nodes/%s", n.Name)
	body, err := JSONReader(n)
	if err != nil {
		return
	}

	err = e.client.magicRequestDecoder("PUT", url, body, &node)
	return
}

// Delete removes a node on the Chef server
//
// Chef API docs: https://docs.chef.io/api_chef_server.html#nodes-name
func (e *NodeService) Delete(name string) (err error) {
	err = e.client.magicRequestDecoder("DELETE", "nodes/"+name, nil, nil)
	return
}
