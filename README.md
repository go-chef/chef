[![Stories in Ready](https://badge.waffle.io/go-chef/chef.png?label=ready&title=Ready)](https://waffle.io/go-chef/chef)
[![Build Status](https://app.wercker.com/status/9cfd4b53ea24e0894904067f283e4cf8/s "wercker status")](https://app.wercker.com/project/bykey/9cfd4b53ea24e0894904067f283e4cf8)
[![Coverage Status](https://coveralls.io/repos/go-chef/chef/badge.png?branch=master)](https://coveralls.io/r/go-chef/chef?branch=master)

# Chef Server API Client Library in Golang
This is a Library that you can use to write tools to interact with the chef server. 

## Install
 
    go get github.com/go-chef/chef

## Test
  
    go get -t github.com/go-chef/chef
    go test -v github.com/go-chef/chef

## SSL

  If you run into an SSL verification problem when trying to connect to a ssl server with self signed certs setup your config object with `SkipSSL: true`

## Usage
This example is setting up a basic client that you can use to interact with all the service endpoints (clients, nodes, cookbooks, etc.)
More usage examples can be found in the [examples](examples) directory.
 
     package main
     
     import (
     	"encoding/json"
     	"fmt"
     	"io/ioutil"
     	"log"
     	"os"
     
     	"github.com/go-chef/chef"
     )
     
     func main() {
     	// read a client key
     	key, err := ioutil.ReadFile("key.pem")
     	if err != nil {
     		fmt.Println("Couldn't read key.pem:", err)
         os.Exit(1)
       }
     
     	// build a client
     	client, err := chef.NewClient(&chef.Config{
     		Name: "foo",
     		Key:  string(key),
     		// goiardi is on port 4545 by default. chef-zero is 8889
     		BaseURL: "http://localhost:4545",
     	})
     	if err != nil {
     		fmt.Println("Issue setting up client:", err)
       }
     
       // List Cookbooks
       cookList := client.Cookbooks.List()
       if err != nil {
         fmt.Println("Issue listing cookbooks:", err)
       }
     
       // Print out the list
       fmt.Println(cookList)
     }

