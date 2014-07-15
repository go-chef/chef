package chef

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// ChefVersion that we pretend to emulate
const ChefVersion = "11.12.0"

// AuthConfig representing a client and a private key used for encryption
type AuthConfig struct {
	privateKey *rsa.PrivateKey
	clientName string
}

// Client is vessel for public methods used against the chef-server
type Client struct {
	Auth    *AuthConfig
	BaseURL *url.URL
	client  *http.Client

	Environments *EnvironmentService
	Nodes        *NodeService
}

// Config contains the configuration options for a chef client
type Config struct {
	Name    string
	Key     string
	BaseURL string
	SkipSSL bool
}

// NewClient is the client generator used to instantiate a client for talking to a chef-server
// It is a simple constructor for the Client struct intended as a easy interface for issuing
// signed requests
func NewClient(cfg *Config) (*Client, error) {
	pk, err := privateKeyFromString([]byte(cfg.Key))
	if err != nil {
		return nil, err
	}

	baseUrl, _ := url.Parse(cfg.BaseURL)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.SkipSSL},
	}

	c := &Client{
		Auth: &AuthConfig{
			privateKey: pk,
			clientName: cfg.Name,
		},
		client:  &http.Client{Transport: tr},
		BaseURL: baseUrl,
	}
	c.Environments = &EnvironmentService{client: *c}
	c.Nodes = &NodeService{client: *c}
	return c, nil
}

// MakeRequest performs a signed request for the chef client
func (c *Client) MakeRequest(method string, requestUrl string, body io.Reader) (*http.Request, error) {
	relativeUrl, err := url.Parse(requestUrl)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(relativeUrl)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	// don't have to check this works, signRequest only emits error when signing hash is not valid, and we baked that in
	c.Auth.SignRequest(req)

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, res.Body)
		} else {
			err = json.NewDecoder(res.Body).Decode(v)
		}
	}

	return res, nil
}

// SignRequest modifies headers of an http.Request
func (ac AuthConfig) SignRequest(request *http.Request) error {
	// sanitize the path for the chef-server
	// chef-server doesn't support '//' in the Hash Path.
	var endpoint string
	if request.URL.Path != "" {
		endpoint = path.Clean(request.URL.Path)
		request.URL.Path = endpoint
	} else {
		endpoint = request.URL.Path
	}

	request.Header.Set("Method", request.Method)
	request.Header.Set("Hashed Path", hashStr(endpoint))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-Chef-Version", ChefVersion)
	request.Header.Set("X-Ops-Timestamp", time.Now().UTC().Format(time.RFC3339))
	request.Header.Set("X-Ops-UserId", ac.clientName)
	request.Header.Set("X-Ops-Sign", "algorithm=sha1;version=1.0")
	request.Header.Set("X-Ops-Content-Hash", calcBodyHash(request))

	// To validate the signature it seems to be very particular
	var content string
	for _, key := range []string{"Method", "Hashed Path", "X-Ops-Content-Hash", "X-Ops-Timestamp", "X-Ops-UserId"} {
		content += fmt.Sprintf("%s:%s\n", key, request.Header.Get(key))
	}
	content = strings.TrimSuffix(content, "\n")
	// generate signed string of headers
	// Since we've gone through additional validation steps above,
	// we shouldn't get an error at this point
	signature, err := generateSignature(ac.privateKey, content)
	if err != nil {
		return err
	}

	// TODO: THIS IS CHEF PROTOCOL SPECIFIC
	// Signature is made up of n 60 length chunks
	base64sig := base64BlockEncode(signature, 60)

	// roll over the auth slice and add the apropriate header
	for index, value := range base64sig {
		request.Header.Set(fmt.Sprintf("X-Ops-Authorization-%d", index+1), string(value))
	}

	return nil
}

// modified from goiardi calcBodyHash
func calcBodyHash(r *http.Request) string {
	var bodyStr string

	if r.Body == nil {
		bodyStr = ""
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		bodyStr = buf.String()
	}

	// Since we're not setting the encoded slice limit
	// we can safely call out [0]
	chkHash := hashStr(bodyStr)
	return chkHash
}

// privateKeyFromString parses an RSA private key from a string
func privateKeyFromString(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, fmt.Errorf("block size invalid for '%s'", string(key))
	}
	rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsaKey, nil
}
