package chef

// PolicyService  is the service for interacting with chef server policies endpoint
type PolicyService struct {
	client *Client
}

// PolicyGetResponse is returned from the chef-server for Get Requests to /policies
type PoliciesGetResponse map[string]Policy

type Policy struct {
	Uri       string                 `json:"uri,omitempty"`
	Revisions map[string]interface{} `json:"revisions,omitempty"`
}

// List lists the policies in the Chef server.
// Chef API docs: https://docs.chef.io/api_chef_server/#policies
func (e *PolicyService) List() (data PoliciesGetResponse, err error) {
	err = e.client.magicRequestDecoder("GET", "policies", nil, &data)
	return
}
