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

// Environment represents the native Go version of the deserialized Environment type
type Environment struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ChefType    string      `json:"chef_type"`
	Attributes  interface{} `json:"attributes,omitempty"`
	// DefaultAttributes  interface{}       `json:"default_attributes,omitempty"`
	// OverrideAttributes interface{}       `json:"override_attributes,omitempty"`
	JsonClass        string            `json:"json_class,omitempty"`
	CookbookVersions map[string]string `json:"cookbook_versions"`
	i                int64             // current reading index
}

func (b *Environment) Read(p []byte) (n int, err error) {
	if b == new(Environment) {
		return 0, nil
	}

	if len(p) == 0 {
		return 0, nil
	}

	buf, err := json.Marshal(&b) // should save this
	if err != nil {
		fmt.Println("error doing json")
		return 0, err
	}

	fmt.Println(fmt.Sprintf("sending: %s", buf))
	if b.i >= int64(len(buf)) {
		b.i = 0
		return 0, io.EOF
	}

	n = copy(p, buf[b.i:])
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
	err = e.client.magicRequestDecoder("PUT", path, environment, nil)
	return
}
