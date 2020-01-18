package chef

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// TODO: Search  doesn't look implemented

const algorithm string = "aes-256-gcm"
const supportedVersion int = 3

var errUnsupportedVersion = errors.New("Only Version 3 of the encrypted databag item format is supported")

// VaultService is the service for interacting with the chef server data endpoint
type VaultService struct {
	client *Client
}

// VaultItem wraps a vault databag and its keys
type VaultItem struct {
	DataBagItem  *DataBagItem
	Keys         *VaultItemKeys
	Name         string
	Vault        string
	VaultService *VaultService
}

// VaultItemKeys contains the client keys associated with a VaultItem
type VaultItemKeys struct {
	DataBagItem *DataBagItem
	Name        string
}

// VaultListResult is the list of vaults returned by chef-api when listing
// http://docs.getchef.com/api_chef_server.html#data
type VaultListResult map[string]string

// String makes VaultListResult implement the string result
func (vlr VaultListResult) String() (out string) {
	for k, v := range vlr {
		out += fmt.Sprintf("%s => %s\n", k, v)
	}
	return out
}

// List returns a list of vaults on the server.  The official implementation of this method loads each databag
// and makes some assumptions based on the number of keys and their names.
// See https://sourcegraph.com/github.com/chef/chef-vault@0d008c3/-/blob/lib/chef/knife/vault_list.rb#L38-51
func (vs *VaultService) List() (*VaultListResult, error) {
	databags, err := vs.client.DataBags.List()
	if err != nil {
		return nil, err
	}

	//	keys := map[string]bool{}

	// The following method doesn't work
	// The _keys entries are items not vaults
	//for name := range *databags {
	//if vaultName := strings.TrimSuffix(name, "_keys"); vaultName != name {
	//keys[vaultName] = true
	//}
	//}

	vaults := VaultListResult{}
	for k, v := range *databags {
		items, err := vs.ListItems(k)
		if err != nil {
			return &vaults, err
		}
		if len(items) > 0 {
			vaults[k] = v
		}
	}

	//for k := range keys {
	//if v, ok := (*databags)[k]; ok {
	//vaults[k] = v
	//}
	//}

	return &vaults, nil
}

// CreateItem creates a vault item and keys databag
func (vs *VaultService) CreateItem(vaultName, itemName string) (*VaultItem, error) {
	sharedSecret := *generateSecret()
	userEncodedSecret, err := EncodeSharedSecret(vs.client.Auth.PrivateKey, sharedSecret)
	if err != nil {
		return nil, err
	}

	keysItem := map[string]interface{}{
		"admins":                  []string{vs.client.Auth.ClientName},
		"clients":                 []string{},
		"id":                      keysItemName(itemName),
		"mode":                    "default",
		"search_query":            []string{},
		vs.client.Auth.ClientName: userEncodedSecret,
	}
	primaryItem := map[string]interface{}{
		"id": itemName,
	}

	// Create the vault's databag unless it exists
	databag := &DataBag{Name: vaultName}
	if _, err := vs.client.DataBags.Create(databag); err != nil {
		if !strings.Contains(err.Error(), " 409") {
			return nil, err
		}
	}

	// Create client keys databag
	if err := vs.client.DataBags.CreateItem(vaultName, keysItem); err != nil {
		return nil, err
	}

	// Create primary databag
	if err := vs.client.DataBags.CreateItem(vaultName, primaryItem); err != nil {
		return nil, err
	}

	databagItem := DataBagItem(primaryItem)
	databagKeysItem := DataBagItem(keysItem)

	return &VaultItem{
		DataBagItem:  &databagItem,
		Name:         itemName,
		Vault:        vaultName,
		VaultService: vs,
		Keys: &VaultItemKeys{
			DataBagItem: &databagKeysItem,
			Name:        keysItemName(itemName),
		},
	}, nil
}

// DeleteItem deletes an item from a data bag
//   Chef API Docs: http://docs.getchef.com/api_chef_server.html#id22
func (vs *VaultService) DeleteItem(vaultName, vaultItem string) error {
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
	err := vs.client.DataBags.DeleteItem(vaultName, keysItemName(vaultItem))
	if re := notFound(err); re != nil {
		return re
	}

	// Delete primary databag item
	err = vs.client.DataBags.DeleteItem(vaultName, vaultItem)
	if re := notFound(err); re != nil {
		return re
	}

	return nil
}

// GetItem fetches a VaultItem and loads client keys
func (vs *VaultService) GetItem(vaultName, itemName string) (*VaultItem, error) {
	databagItem, err := vs.client.DataBags.GetItem(vaultName, itemName)
	if err != nil {
		return nil, err
	}

	item := &VaultItem{
		Name:         itemName,
		VaultService: vs,
		DataBagItem:  &databagItem,
		Vault:        vaultName,
	}

	vaultKeys, err := item.loadKeys(vs.client, vaultName, itemName)
	if err != nil {
		return nil, err
	}

	item.Keys = vaultKeys

	return item, nil
}

// ListItems lists the items in a vault
func (vs *VaultService) ListItems(vault string) ([]string, error) {
	bagitems, err := vs.client.DataBags.ListItems(vault)
	if err != nil {
		return nil, err
	}
	itemnames := onlyitemnames(bagitems)
	return itemnames, err
}

func onlyitemnames(items *DataBagListResult) []string {
	var itemnames []string
	for k, _ := range *items {
		// if we have and item and item_keys assume we have a vault item
		if _, ok := (*items)[k+"_keys"]; !ok {
			continue
		}
		itemnames = append(itemnames, k)
	}
	return itemnames
}

// UpdateItem sets the item data, encrypts with a shared key, and then encrypts the shared key with each authorized client key in the <item>_keys data bag
func (vs *VaultService) UpdateItem(item *VaultItem, data map[string]interface{}) error {
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

	err = vs.client.DataBags.UpdateItem(item.Vault, item.Name, itemData)
	if err != nil {
		return err
	}

	databagData := DataBagItem(itemData)

	item.DataBagItem = &databagData

	return nil
}

// ListItemAdmins return a list of admin users for the item
func (vs *VaultService) ListItemAdmins(item *VaultItem) []string {
	return (*(*(*item).Keys).DataBagItem).(map[string]interface{})["admins"].([]string)
}

// ListItemClients return a list of client users for the item
func (vs *VaultService) ListItemClients(item *VaultItem) []string {
	return (*(*(*item).Keys).DataBagItem).(map[string]interface{})["clients"].([]string)
}

// UpdateItemAdmins sets the list of admin users for the item
func (vs *VaultService) UpdateItemAdmins(item *VaultItem, admins []string) {
	(*(*(*item).Keys).DataBagItem).(map[string]interface{})["admins"] = admins
	return
}

// UpdateItemClients sets the list of clients for the item
func (vs *VaultService) UpdateItemClients(item *VaultItem, clients []string) {
	(*(*(*item).Keys).DataBagItem).(map[string]interface{})["clients"] = clients
	return
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
		"version":        float64(supportedVersion),
		"cipher":         algorithm,
	}, nil
}
