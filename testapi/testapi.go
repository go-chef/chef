// Test the go-chef/chef chef server api /organizations endpoints against a live server
package testapi

import (
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strconv"

	chef "github.com/go-chef/chef"
)

// client exercise the chef server api
func Client() *chef.Client {
	// Pass in the database and chef-server api credentials.
	user := os.Args[2]
	keyfile := os.Args[3]
	chefurl := os.Args[4]
	skipssl, err := strconv.ParseBool(os.Args[5])
	if err != nil {
		skipssl = true
	}
	version := "1.0"
	if len(os.Args) > 6 {
		version = os.Args[6]
	}

	// Create a client for access
	return buildClient(user, keyfile, chefurl, skipssl, version)
}

// buildClient creates a connection to a chef server using the chef api.
// goiardi uses port 4545 by default, chef-zero uses 8889, chef-server uses 443
func buildClient(user string, keyfile string, baseurl string, skipssl bool, version string) *chef.Client {
	key := clientKey(keyfile)
	client, err := chef.NewClient(&chef.Config{
		Name:                  user,
		Key:                   string(key),
		BaseURL:               baseurl,
		SkipSSL:               skipssl,
		RootCAs:               chefCerts(),
		AuthenticationVersion: version,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue setting up client:", err)
		os.Exit(1)
	}
	return client
}

// clientKey reads the pem file containing the credentials needed to use the chef client.
func clientKey(filepath string) string {
	key, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read key.pem: %+v, %+v", filepath, err)
		os.Exit(1)
	}
	return string(key)
}

// chefCerts creats a cert pool for the self signed certs
// reference https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
func chefCerts() *x509.CertPool {
	const localCertFile = "/var/opt/opscode/nginx/ca/testhost.crt"
	certPool, _ := x509.SystemCertPool()
	if certPool == nil {
		certPool = x509.NewCertPool()
	}
	// Read in the cert file
	certs, err := os.ReadFile(localCertFile)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", localCertFile, err)
	}
	// Append our cert to the system pool
	if ok := certPool.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}
	return certPool
}
