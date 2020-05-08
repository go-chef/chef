package chef

// PolicyGroupService  is the service for interacting with chef server policies endpoint
type PolicyGroupService struct {
	client *Client
}

// PolicyGroupGetResponse is returned from the chef-server for Get Requests to /policy_groups
type PolicyGroupGetResponse map[string]PolicyGroup

type PolicyGroup struct {
	Uri      string              `json:"uri,omitempty"`
	Policies map[string]Revision `json:"policies,omitempty"`
}

type Revision map[string]string

// List lists the policy groups in the Chef server.
// Chef API docs: https://docs.chef.io/api_chef_server/#policy_groups
func (e *PolicyGroupService) List() (data PolicyGroupGetResponse, err error) {
	err = e.client.magicRequestDecoder("GET", "policy_groups", nil, &data)
	return
}

// Delete deletes a policy group.
// DELETE /policy_groups/GROUP
// Chef API docs: https://docs.chef.io/api_chef_server/#policy_groups
func (e *PolicyGroupService) Delete(policyGroupName string) (data PolicyGroupGetResponse, err error) {
	err = e.client.magicRequestDecoder("DELETE", "policy_groups/" + policyGroupName, nil, &data)
	return
}


// policy_group oc_chef_wm_policy_groups.erl  
  // GET
// policy_group/GN oc_chef_wm_named_policy_group.erl
 // DELETE, GET
// policy_group/GN/policies/PN  not sure of the path, could be revisons instead of policies oc_chef_wm_named_policy_named_revision.erl
 // DELETE, GET, PUT
