package chef

import "fmt"

// PolicyService  is the service for interacting with chef server policy endpoint
type PolicyService struct {
	client *Client
}

// CookbookItem represents a object of cookbook file data
type PolicyItem struct {
	Url         string `json:"url,omitempty"`
	Path        string `json:"path,omitempty"`
	Name        string `json:"name,omitempty"`
	Checksum    string `json:"checksum,omitempty"`
	Specificity string `json:"specificity,omitempty"`
}

// PolicyListResult is the summary info returned by chef-api when listing
// http://docs.opscode.com/api_chef_server.html#cookbooks
type PolicyListResult map[string]PolicyVersions

// CookbookVersions is the data container returned from the chef server when listing all cookbooks
type PolicyVersions struct {
	Url      string            `json:"url,omitempty"`
	Versions []CookbookVersion `json:"versions,omitempty"`
}

// String makes PolicyListResult implement the string result
func (p PolicyListResult) String() (out string) {
	for k, v := range p {
		out += fmt.Sprintf("%s => %s\n", k, v.Url)
		for _, i := range v.Versions {
			out += fmt.Sprintf(" * %s\n", i.Version)
		}
	}
	return out
}

// Get retruns a CookbookVersion for a specific cookbook
//  GET /cookbooks/name
func (p *PolicyService) Get(name string, organization string) (data CookbookVersion, err error) {
	path := fmt.Sprintf("organizations/%s/policies", organization)
	err = p.client.magicRequestDecoder("GET", path, nil, &data)
	return
}

// ListVersions lists the cookbooks available on the server limited to numVersions
//   Chef API docs: https://docs.chef.io/api_chef_server.html#cookbooks-name
func (p *PolicyService) ListAvailableQuery(path string) (data PolicyListResult, err error) {
	err = p.client.magicRequestDecoder("GET", path, nil, &data)
	return
}

// List returns a CookbookListResult with the latest versions of cookbooks available on the server
func (p *PolicyService) ListByPolicyGroup(group string) (PolicyListResult, error) {
	path := fmt.Sprintf("policy_groups/%s", group)
	return p.ListAvailableQuery(path)
}

func (p *PolicyService) ListAllPolicyGroups() (PolicyListResult, error) {
	return p.ListAvailableQuery("policy_groups")
}
