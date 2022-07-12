//go:build httpvar
// +build httpvar

package chef

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
)

// TestNewClientProxy2
// Verify that getting proxy information from environment variables works
// This test needs to be run seperately from other tests
// http.ProxyFromEnvironment will only get the proxy value once, the first time it is called
func TestNewClientProxy2(t *testing.T) {
	//test proxy from environment variable
	os.Setenv("https_proxy", "https://8.8.8.8:8000")
	cfg := &Config{Name: "testclient", Key: privateKeyPKCS1, SkipSSL: false, Timeout: 1}
	chefClient, err := NewClient(cfg)
	assert.Nil(t, err, "Create client")
	request, err := chefClient.NewRequest("GET", "https://test.com", nil)
	assert.Nil(t, err, "Create request")

	eurl := &url.URL{Scheme: "https", Host: "8.8.8.8:8000"}
	trurl, err := chefClient.client.Transport.(*http.Transport).Proxy(request)
	assert.Equal(t, *eurl, *trurl, "proxy value from environment variable")

	tr := chefClient.client.Transport.(*http.Transport)
	assert.Equal(t, reflect.ValueOf(tr.Proxy).Pointer(),
		reflect.ValueOf(http.ProxyFromEnvironment).Pointer(),
		"Proxy set from http proxyfromenvironment function")
}
