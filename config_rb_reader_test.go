package chef

import (
	"strings"
	"testing"
)

func TestNewClientRb(t *testing.T) {
	clientRegistry["client_key"] = func(s []string, path string, m *ConfigRb) {
		str := StringParserForMeta(s)
		data := strings.Split(str, "/")
		m.ClientKey = data[len(data)-1]
	}
	data := `current_dir = File.dirname(__FILE__)
		log_level                :info
		log_location             STDOUT
		node_name                "test"
		client_key               "#{current_dir}/test.pem"
		chef_server_url          "https://server/organizations/test"
		cookbook_path            ["#{current_dir}/../cookbooks"]`
	cb, err := NewClientRb(data, "")
	if err != nil {
		t.Error("unable to read config.rb file")
	}
	if cb.ClientKey != "test.pem" {
		t.Error("client.pem have invalid path")
	}
	if cb.ChefServerUrl != "https://server/organizations/test" {
		t.Error("invalid chef server url read from config.rb")
	}
	if cb.NodeName != "test" {
		t.Error("invalid node name read from config.rb")

	}
}
