package chef

import (
	"encoding/json"
	"fmt"
	"io"
)

type RoleService struct {
	client *Client
}

type RoleListResult map[string]string

// Role represents the native Go version of the deserialized Role type
type Role struct {
	Name               string      `json:"name"`
	ChefType           string      `json:"chef_type"`
	Description        string      `json:"description"`
	RunList            RunList     `json:"run_list"`
	DefaultAttributes  interface{} `json:"default_attributes,omitempty"`
	OverrideAttributes interface{} `json:"override_attributes,omitempty"`
	JsonClass          string      `json:"json_class,omitempty"`
	i                  int64       // current reading index
	buf                []byte
}

func (b *Role) Read(p []byte) (n int, err error) {
	if b == new(Role) {
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

// String makes RoleListResult implement the string result
func (e RoleListResult) String() (out string) {
	for k, v := range e {
		out += fmt.Sprintf("%s => %s\n", k, v)
	}
	return out
}

// List lists the roles in the Chef server.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id31
func (e *RoleService) List() (data *RoleListResult, err error) {
	err = e.client.magicRequestDecoder("GET", "roles", nil, &data)
	return
}

// Create a new role in the Chef server.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id32
func (e *RoleService) Create(role *Role) (err error) {
	path := fmt.Sprintf("roles")
	err = e.client.magicRequestDecoder("POST", path, role, nil)
	return
}

// Delete a role from the Chef server.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id33

// Get gets a role from the Chef server.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id34
func (e *RoleService) Get(name string) (data *Role, err error) {
	path := fmt.Sprintf("roles/%s", name)
	err = e.client.magicRequestDecoder("GET", path, nil, &data)
	return
}

// Update a role in the Chef server.
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id35
func (e *RoleService) Put(role *Role) (err error) {
	path := fmt.Sprintf("roles/%s", role.Name)
	err = e.client.magicRequestDecoder("PUT", path, role, nil)
	return
}

// Get a list of environments have have environment specific run-lists for the given role
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id36

// Get the environment-specific run-list for  a role
//
// Chef API docs: http://docs.getchef.com/api_chef_server.html#id37
