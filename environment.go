package chef

import "fmt"

// Environment has a Reader, hey presto
type EnvironmentService struct {
	client Client
}

// Environment represents the native Go version of the deserialized Environment type
type Environment struct {
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	ChefType           string      `json:"chef_type"`
	DefaultAttributes  interface{} `json:"default_attributes"`
	OverrideAttributes interface{} `json:"override_attributes"`
	JsonClass          string      `json:"chef_environment"`
	CookbookVersions   interface{} `json:"cookbook_versions"`
}

// List lists the environments in the Chef server.
//
// Chef API docs: http://docs.opscode.com/api_chef_server.html#id14
func (e *EnvironmentService) List() (data map[string]string, err error) {
	req, err := e.client.MakeRequest("GET", "environments", nil)
	if err != nil {
		return
	}

	_, err = e.client.Do(req, &data)
	if err != nil {
		return
	}

	return
}

// Get gets an environment from the Chef server.
//
// Chef API docs: http://docs.opscode.com/api_chef_server.html#id17
func (e *EnvironmentService) Get(name string) (*Environment, error) {
	url := fmt.Sprintf("environments/%s", name)
	req, err := e.client.MakeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	environment := new(Environment)
	_, err = e.client.Do(req, &environment)
	if err != nil {
		return nil, err
	}

	return environment, nil
}
