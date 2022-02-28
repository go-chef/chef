package chef

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ConfigRb struct {
	ClientKey     string
	ChefServerUrl string
	NodeName      string
}

type clientFunc func(s []string, path string, m *ConfigRb)

var clientRegistry map[string]clientFunc

func init() {
	clientRegistry = make(map[string]clientFunc, 2)
	clientRegistry["client_key"] = configKeyParser
	clientRegistry["chef_server_url"] = configServerParser
	clientRegistry["node_name"] = configNodeNameParser

}
func NewClientRb(data, path string) (c ConfigRb, err error) {
	linesData := strings.Split(data, "\n")
	if len(linesData) < 3 {
		return c, errors.New("not much info")
	}
	for _, i := range linesData {
		key, value := getKeyValue(strings.TrimSpace(i))
		if fn, ok := clientRegistry[key]; ok {
			fn(value, path, &c)
		}
	}
	return c, err
}
func configKeyParser(s []string, path string, c *ConfigRb) {
	str := StringParserForMeta(s)
	data := strings.Split(str, "/")
	size := len(data)
	if size > 0 {
		keyPath := filepath.Join(path, data[size-1])
		keyData, err := ioutil.ReadFile(keyPath)
		if err != nil {
			fmt.Println("error in reading pem file at: ", keyPath)
			os.Exit(1)
		}
		c.ClientKey = string(keyData)
	}
}
func configServerParser(s []string, path string, c *ConfigRb) {
	c.ChefServerUrl = StringParserForMeta(s)
}
func configNodeNameParser(s []string, path string, c *ConfigRb) {
	c.NodeName = StringParserForMeta(s)
}
