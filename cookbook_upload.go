package chef

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MaxCookbookUploadRetries = 5
const MaxCookbookUploadRetrySleep = 2 * time.Second

func (c *CookbookService) UploadFromForce(path string) (*Cookbook, error) {
	return c.UploadFrom(path, true, true)
}

func (c *CookbookService) UploadFromFrozen(path string) (*Cookbook, error) {
	return c.UploadFrom(path, true, false)
}

// UploadFrom uploads a cookbook from a given path to the Chef server
// frozen indicates whether the cookbook should be frozen, and force indicates
// whether the cookbook should be force uploaded (override an existing frozen version)
//
// * Parses local cookbook.
// * Creates a sandbox with the cookbook's item checksums.
// * Uploads each missing cookbook item into the sandbox.
// * Commits the sandbox.
// * Uploads the cookbook version manifest.
func (c *CookbookService) UploadFrom(path string, frozen bool, force bool) (*Cookbook, error) {
	cookbook, err := NewCookbookFromPath(path)
	if err != nil {
		return nil, err
	}
	cookbook.Frozen = frozen

	metadataJsonExists := false
	cookbookItems := cookbook.AllItemsByChecksum()
	checksums := []string{}
	for k, v := range cookbookItems {
		checksums = append(checksums, k)
		if v.Path == "metadata.json" {
			metadataJsonExists = true
		}
	}

	// If the metadata JSON does not exist, generate it from the metadata
	if !metadataJsonExists {
		metadataJsonBytes, err := json.MarshalIndent(cookbook.Metadata, "", "    ")
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(filepath.Join(path, "metadata.json"), metadataJsonBytes, os.ModePerm)
		if err != nil {
			return nil, err
		}
		// Ensure the generated metadata.json is removed
		defer os.Remove(filepath.Join(path, "metadata.json"))

		metadataJsonChecksum, err := fileMD5Checksum(filepath.Join(path, "metadata.json"))
		if err != nil {
			return nil, err
		}

		metadataJsonCi := CookbookItem{
			Name:        "metadata.json",
			Path:        "metadata.json",
			Checksum:    metadataJsonChecksum,
			Specificity: "default",
		}

		// Update associated data
		cookbook.RootFiles = append(cookbook.RootFiles, metadataJsonCi)
		cookbookItems[metadataJsonChecksum] = metadataJsonCi
		checksums = append(checksums, metadataJsonChecksum)
	}

	// Create the new sandbox
	sandboxPostResp, err := c.client.Sandboxes.Post(checksums)
	if err != nil {
		return nil, err
	}

	sandboxId := sandboxPostResp.ID
	sandboxChecksums := sandboxPostResp.Checksums

	// Upload cookbook files to sandbox
	for checksum, checksumDetails := range sandboxChecksums {
		// Skip files that do not require upload
		if !checksumDetails.Upload {
			continue
		}

		// TODO: Parallelize the uploads
		cookbookItem := cookbookItems[checksum]
		itemPath := filepath.Join(path, cookbookItem.Path)
		err = c.UploadCookbookItem(itemPath, &cookbookItem, checksumDetails.Url)
		if err != nil {
			return nil, err
		}
	}

	// Commit the sandbox
	//
	// Retries are performed to reflect the Ruby implementation
	// (the upload target may not yet have replicated the uploaded data)
	retries := 0

	for retries < MaxCookbookUploadRetries {
		_, err = c.client.Sandboxes.Put(sandboxId)
		if err != nil {
			// Retry 400 errors until max retries
			if strings.Contains(err.Error(), ": 400") {
				time.Sleep(MaxCookbookUploadRetrySleep)
				retries++
				continue
			}

			return nil, err
		} else {
			break
		}
	}

	if retries >= MaxCookbookUploadRetries {
		return nil, err
	}

	err = c.UploadVersion(cookbook, force)
	if err != nil {
		return nil, err
	}

	return cookbook, nil
}

// UploadVersion puts a specific version of a cookbooks to the server api
// If force is true, ?force=true is appended to the cookbook URL.
//
//	PUT /cookbook/foo/1.2.3
//	Chef API docs: https://docs.chef.io/server/api_chef_server/#put-6
func (c *CookbookService) UploadVersion(cookbook *Cookbook, force bool) error {
	url := fmt.Sprintf("cookbooks/%s/%s", cookbook.CookbookName, cookbook.Version)

	if force {
		url = fmt.Sprintf("%s?force=true", url)
	}

	manifest, err := cookbook.ManifestJsonForApi(c.client.Auth.ServerApiVersion)
	if err != nil {
		return err
	}

	req, err := c.client.NewRequest("PUT", url, manifest)
	if err != nil {
		return err
	}

	_, err = c.client.Do(req, nil)

	return err
}

// Uploads a single cookbook item to a given URL.
// Requires the cookbook item's local path, the CookbookItem, and the target URL.
func (c *CookbookService) UploadCookbookItem(itemPath string, item *CookbookItem, url string) error {
	md5B64Checksum, err := md5Base64Checksum(item.Checksum)
	if err != nil {
		return err
	}

	itemFile, err := os.Open(itemPath)
	if err != nil {
		return err
	}
	defer itemFile.Close()

	uploadReq, err := c.client.NewRequest("PUT", url, itemFile)
	if err != nil {
		return err
	}

	// Get file length for content-length header (required to prevent the Go http client from chunking the request)
	fileInfo, err := os.Stat(itemPath)
	if err != nil {
		return err
	}

	// Set minimum headers for file upload
	uploadReq.ContentLength = fileInfo.Size()
	uploadReq.Header.Set("content-md5", md5B64Checksum)
	uploadReq.Header.Set("content-type", "application/x-binary")
	uploadReq.Header.Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	_, err = c.client.Do(uploadReq, nil)

	return err
}
