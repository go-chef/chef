package chef

type UpdatedSinceService struct {
	client *Client
}

// UpdatedSince represents the body of the returned information.
type UpdatedSince {
	Action string
	Id string
	Path string
}

// Since gets available cookbook version information.
//
// https://docs.chef.io/api_chef_server/#updated_since
// client should have a base url with a specified organization
func (e UpdatedSinceService) Get(sequenceId string) (updated []UpdatedSince, err error) {
	url := "updated_since"
        if len(filters) > 0 {
                url += "?seq=" + sequenceId
        }
	err = e.client.magicRequestDecoder("GET", url, nil, &updated)
	return
}
