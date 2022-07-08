package chef

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientRb(t *testing.T) {
	clientRegistry["client_key"] = func(s []string, path string, m *ConfigRb) error {
		str := StringParserForMeta(s)
		data := strings.Split(str, "/")
		m.ClientKey = data[len(data)-1]
		return nil
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
	assert.Equal(t, "test.pem", cb.ClientKey)
	assert.Equal(t, "https://server/organizations/test", cb.ChefServerUrl)
	assert.Equal(t, "test", cb.NodeName)
}
