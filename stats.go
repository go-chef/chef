package chef

type StatsService struct {
	client *Client
}

// Stats: represents the body of the returned information.
type Stats []map[string]interface{}

// Stat gets the frontend & backend server information.
//
// https://docs.chef.io/api_chef_server/
//
// ?format = text or json or nothing, text is supposed to work but returns a 406
func (e *StatsService) Get(format string, user string, password string) (data Stats, err error) {
	if format == "" {
		format = "json"
	}
	err = e.client.basicRequestDecoder("GET", "_stats?format="+format, nil, &data, user, password)
	return
}
