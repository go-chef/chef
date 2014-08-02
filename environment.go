package chef

import (
	"encoding/json"
	"fmt"
	"io"
)

// Environment has a Reader, hey presto
type EnvironmentService struct {
	client *Client
}

type EnvironmentListResult map[string]string
type EnvironmentCreateResult map[string]string

// Environment represents the native Go version of the deserialized Environment type
type Environment struct {
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	ChefType           string            `json:"chef_type"`
	Attributes         interface{}       `json:"attributes,omitempty"`
	DefaultAttributes  interface{}       `json:"default_attributes,omitempty"`
	OverrideAttributes interface{}       `json:"override_attributes,omitempty"`
	JsonClass          string            `json:"json_class,omitempty"`
	CookbookVersions   map[string]string `json:"cookbook_versions"`
	i                  int64             // current reading index
	buf                []byte
}

func (b *Environment) Read(p []byte) (n int, err error) {
	if b == new(Environment) {
		return 0, nil
	}

	if len(p) == 0 {
		return 0, nil
	}

	if b.buf == nil {
		b.buf, err = json.Marshal(&b)
		if err != nil {
			return 0, err
		}
	}

	if b.i >= int64(len(b.buf)) {
		b.i = 0
		b.buf = nil
		return 0, io.EOF
	}

	n = copy(p, b.buf[b.i:])
	b.i += int64(n)
	return
}

// String makes EnvironmentListResult implement the string result
func (e EnvironmentListResult) String() (out string) {
	for k, v := range e {
		out += fmt.Sprintf("%s => %s\n", k, v)
	}
	return out
}

// String makes EnvironmentCreateResult implement the string result
func (e EnvironmentCreateResult) String() (out string) {
	for k, v := range e {
		out += fmt.Sprintf("%s => %s\n", k, v)
	}
	return out
}

// List lists the environments in the Chef server.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id14
func (e *EnvironmentService) List() (data *EnvironmentListResult, err error) {
	err = e.client.magicRequestDecoder("GET", "environments", nil, &data)
	return
}

// Create an environment in the Chef server.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id15
func (e *EnvironmentService) Create(environment *Environment) (data *EnvironmentCreateResult, err error) {
	err = e.client.magicRequestDecoder("POST", "environments", environment, &data)
	return
}

// Delete an environment from the Chef server.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id16

// Get gets an environment from the Chef server.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id17
func (e *EnvironmentService) Get(name string) (data *Environment, err error) {
	path := fmt.Sprintf("environments/%s", name)
	err = e.client.magicRequestDecoder("GET", path, nil, &data)
	return
}

// Write an environment to the Chef server.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id18
func (e *EnvironmentService) Put(environment *Environment) (data *Environment, err error) {
	path := fmt.Sprintf("environments/%s", environment.Name)
	err = e.client.magicRequestDecoder("PUT", path, environment, &data)
	return
}

// Get the versions of a cookbook for this environment from the Chef server.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id19

// Get a hash of cookbooks and cookbook versions (including all dependencies) that
// are required by the run_list array. Version constraints may be specified using
// the @ symbol after the cookbook name as a delimiter. Version constraints may also
// be present when the cookbook_versions attributes is specified for an environment
// or when dependencies are specified by a cookbook.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id20

// Get a list of cookbooks and cookbook versions that are available to the specified environment.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id21
