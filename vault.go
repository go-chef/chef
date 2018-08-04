package chef

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

const algorithm string = "aes-256-gcm"
const supportedVersion int = 3

var errUnsupportedVersion = errors.New("Only Version 3 encrypted values are supported")

// VaultService is the service for interacting with the chef server data endpoint
type VaultService struct {
	client *Client
}

// Vault is an encrypted data bag
type Vault struct {
	DataBag
}

// VaultItem wraps a vault databag and it's keys
type VaultItem struct {
	DataBagItem  *DataBagItem
	Name         string
	Keys         *VaultItemKeys
	Secret       string
	Vault        string
	VaultService *VaultService
}

// VaultItemKeys contains the client keys associated with a
type VaultItemKeys struct {
	DataBagItem *DataBagItem
	Name        string
}

// VaultListResult is the list of data bags returned by chef-api when listing
// http://docs.getchef.com/api_chef_server.html#data
type VaultListResult map[string]string

// String makes VaultListResult implement the string result
func (d VaultListResult) String() (out string) {
	for k, v := range d {
		out += fmt.Sprintf("%s => %s\n", k, v)
	}
	return out
}

// List returns a list of databags on the server
//   Chef API Docs: http://docs.getchef.com/api_chef_server.html#id18
func (d *VaultService) List() (data *VaultListResult, err error) {
	path := fmt.Sprintf("data")
	err = d.client.magicRequestDecoder("GET", path, nil, &data)
	return
}

// Create adds a data bag to the server
//   Chef API Docs: http://docs.getchef.com/api_chef_server.html#id19
func (d *VaultService) Create(databag *DataBag) (result *DataBagCreateResult, err error) {
	body, err := JSONReader(databag)
	if err != nil {
		return
	}

	err = d.client.magicRequestDecoder("POST", "data", body, &result)
	return
}

// Delete removes a data bag from the server
//   Chef API Docs: ????????????????
func (d *VaultService) Delete(name string) (result *DataBag, err error) {
	path := fmt.Sprintf("data/%s", name)
	err = d.client.magicRequestDecoder("DELETE", path, nil, &result)
	return
}

// ListItems gets a list of the items in a data bag.
//   Chef API Docs: http://docs.getchef.com/api_chef_server.html#id20
func (d *VaultService) ListItems(name string) (data *DataBagListResult, err error) {
	path := fmt.Sprintf("data/%s", name)
	err = d.client.magicRequestDecoder("GET", path, nil, &data)
	return
}

// CreateItem creates a vault item and keys databag
func (d *VaultService) CreateItem(vaultName, vaultItem string) (*VaultItem, error) {
	sharedSecret := *generateSecret()
	userEncodedSecret, err := EncodeSharedSecret(d.client.Auth.PrivateKey, sharedSecret)
	if err != nil {
		return nil, err
	}

	keysItem := map[string]interface{}{
		"admins":                 []string{d.client.Auth.ClientName},
		"clients":                []string{},
		"id":                     keysItemName(vaultItem),
		"mode":                   "default",
		"search_query":           []string{},
		d.client.Auth.ClientName: userEncodedSecret,
	}
	primaryItem := map[string]interface{}{
		"id": vaultItem,
	}

	// Create client keys databag
	if err := d.client.DataBags.CreateItem(vaultName, keysItem); err != nil {
		return nil, err
	}

	// Create primary databag
	if err := d.client.DataBags.CreateItem(vaultName, primaryItem); err != nil {
		return nil, err
	}

	return &VaultItem{}, nil
}

// DeleteItem deletes an item from a data bag
//   Chef API Docs: http://docs.getchef.com/api_chef_server.html#id22
func (d *VaultService) DeleteItem(vaultName, vaultItem string) error {
	notFound := func(err error) error {
		if err != nil {
			if errRes, ok := err.(*ErrorResponse); ok {
				if errRes.Response.StatusCode != 404 {
					return err
				}
			} else {
				return err
			}
		}
		return nil
	}

	// Delete item keys databag item
	err := d.client.DataBags.DeleteItem(vaultName, keysItemName(vaultItem))
	if re := notFound(err); re != nil {
		return re
	}

	// Delete primary databag item
	err = d.client.DataBags.DeleteItem(vaultName, vaultItem)
	if re := notFound(err); re != nil {
		return re
	}

	return nil
}

// GetItem fetches a VaultItem and loads client keys
func (d *VaultService) GetItem(vaultName, itemName string) (*VaultItem, error) {
	databagItem, err := d.client.DataBags.GetItem(vaultName, itemName)
	if err != nil {
		return nil, err
	}

	item := &VaultItem{
		Name:         itemName,
		VaultService: d,
		DataBagItem:  &databagItem,
		Vault:        vaultName,
	}

	vaultKeys, err := item.loadKeys(d.client, vaultName, itemName)
	if err != nil {
		return nil, err
	}

	item.Keys = vaultKeys

	return item, nil
}

// UpdateItem sets the item data, encrypts with a shared key, and then encrypts the shared key with each authorized client key in the <item>_keys data bag
func (d *VaultService) UpdateItem(item *VaultItem, data map[string]interface{}) error {
	itemData := map[string]interface{}{}
	sharedSecret, err := item.sharedSecret()
	if err != nil {
		return err
	}

	for key, value := range data {
		if key == "id" {
			// do nothing
		} else {
			wrappedData, err := json.Marshal(map[string]interface{}{"json_wrapper": value})
			if err != nil {
				return err
			}
			encryptedValue, err := encryptItemValue(wrappedData, sharedSecret)
			itemData[key] = encryptedValue
		}
	}
	itemData["id"] = item.Name

	err = d.client.DataBags.UpdateItem(item.Vault, item.Name, itemData)
	if err != nil {
		return err
	}

	return nil
}

func (i *VaultItem) sharedSecret() ([]byte, error) {
	keys := i.Keys
	auth := i.VaultService.client.Auth

	// Find key entry matching client name
	keyRegistry := (*keys.DataBagItem).(map[string]interface{})
	encryptedSharedSecret, found := keyRegistry[auth.ClientName]
	if !found {
		return nil, errors.New(i.Name + " is not encrypted with your key!")
	}

	rawSecret, err := DecodeSharedSecret(auth.PrivateKey, encryptedSharedSecret.(string))
	if err != nil {
		return nil, err
	}

	encodedSecretArray := sha256.Sum256(rawSecret)
	return encodedSecretArray[:], nil
}

// Decrypt uses the client's PEM to decrypt the shared secret in the keys data bag, and then decrypt the item using the shared secret
func (i *VaultItem) Decrypt() (*map[string]interface{}, error) {
	sharedSecret, err := i.sharedSecret()
	if err != nil {
		return nil, err
	}

	encryptedData := (*i.DataBagItem).(map[string]interface{})
	decryptedItem := make(map[string]interface{})

	// Decrypt each value using the shared secret
	// Each value has the form of "{\"json_wrapper\":{\"foo\":\"b\",\"bar\":[\"b\",\"f\"]}}"
	for key, encryptedValue := range encryptedData {
		if key == "id" {
			decryptedItem[key] = encryptedValue
		} else {
			hash := encryptedValue.(map[string]interface{})
			plaintextValue, err := decryptItemValue(hash, sharedSecret)
			if err != nil {
				return nil, err
			}

			var decodedValue map[string]interface{}
			json.Unmarshal(plaintextValue, &decodedValue)
			decryptedItem[key] = decodedValue["json_wrapper"]
		}
	}

	return &decryptedItem, nil
}

func (i *VaultItem) loadKeys(client *Client, vault, item string) (*VaultItemKeys, error) {
	itemKeysDataBag, err := client.DataBags.GetItem(vault, keysItemName(item))

	if err != nil {
		return nil, err
	}

	return &VaultItemKeys{DataBagItem: &itemKeysDataBag, Name: item}, nil
}

func keysItemName(item string) string {
	return item + "_keys"
}

// Version 3: only
func decryptItemValue(encryptedValue map[string]interface{}, sharedSecret []byte) ([]byte, error) {
	if v, ok := encryptedValue["version"].(float64); !ok || (v != float64(supportedVersion)) {
		return nil, errUnsupportedVersion
	}
	encryptedData := encryptedValue["encrypted_data"].(string)
	encodedIV := encryptedValue["iv"].(string)
	encodedTag := encryptedValue["auth_tag"].(string)

	authTag, err := base64.StdEncoding.DecodeString(encodedTag)
	if err != nil {
		return nil, err
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	iv, err := base64.StdEncoding.DecodeString(encodedIV)
	if err != nil {
		return nil, err
	}

	return DecryptValue(sharedSecret, iv[:], authTag, encryptedBytes)
}

func encryptItemValue(plaintext, secret []byte) (map[string]interface{}, error) {
	nonce, err := iv()
	if err != nil {
		return nil, err
	}

	authTag, encryptedBytes, err := EncryptValue(secret, nonce[:], plaintext)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"encrypted_data": base64.StdEncoding.EncodeToString(encryptedBytes),
		"iv":             base64.StdEncoding.EncodeToString(nonce[:]),
		"auth_tag":       base64.StdEncoding.EncodeToString(authTag),
		"version":        supportedVersion,
		"cipher":         algorithm,
	}, nil
}
