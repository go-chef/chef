package chef

import (
	"fmt"
	"os"
)

// Environment has a Reader, hey presto
type EnvironmentService struct {
	client *Client
}

type EnvironmentListResult map[string]string

// Environment represents the native Go version of the deserialized Environment type
type Environment struct {
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	ChefType           string            `json:"chef_type"`
	DefaultAttributes  interface{}       `json:"default_attributes"`
	OverrideAttributes interface{}       `json:"override_attributes"`
	JsonClass          string            `json:"json_class"`
	CookbookVersions   map[string]string `json:"cookbook_versions"`
}

// String makes EnvironmentListResult implement the string result
func (e EnvironmentListResult) String() (out string) {
	for k, v := range e {
		out += fmt.Sprintf("%s => %s\n", k, v)
	}
	return out
}

// List lists the environments in the Chef server.
//
// Chef API docs: http://docs.opscode.com/api_chef_server.html#id14
func (e *EnvironmentService) List() (data *EnvironmentListResult, err error) {
	err = e.client.magicRequestDecoder("GET", "environments", nil, &data)
	return
}

// Get gets an environment from the Chef server.
//
// Chef API docs: http://docs.opscode.com/api_chef_server.html#id17
func (e *EnvironmentService) Get(name string) (data *Environment, err error) {
	path := fmt.Sprintf("environments/%s", name)
	err = e.client.magicRequestDecoder("GET", path, nil, &data)
	return
}

// Write an environment to the Chef server.
//
// Chef API docs: http://docs.opscode.com/api_chef_server.html#id18
func (e *EnvironmentService) Put(environment *Environment) (err error) {
	path := fmt.Sprintf("environments/%s", environment.Name)
	err = e.client.magicRequestDecoder("PUT", path, environment, os.Stdout)
	return
}
