package testapi

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"os"
)

//random_data makes random byte slice for building junk sandbox data
func random_data(size int) (b []byte) {
	b = make([]byte, size)
	rand.Read(b)
	return
}

// sandbox exercise the chef api
func Sandbox() {
	client := Client()

	// create junk files and sums
	files := make(map[string][]byte)
	sums := make([]string, 10)
	for i := 0; i < 10; i++ {
		data := random_data(128)
		hashstr := fmt.Sprintf("%x", md5.Sum(data))
		files[hashstr] = data
		sums[i] = hashstr
	}

	// TODO: Find a sandbox delete method

	// post the new sums and get a new sandbox id
	postResp, err := client.Sandboxes.Post(sums)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue making request: %+v\n", err)
	}
	fmt.Printf("Create sandboxes %+v\n", postResp)

	// Let's upload the files that postRep thinks we should upload
	for hash, item := range postResp.Checksums {
		if item.Upload == true {
			if hash == "" {
				continue
			}
			// If you were writing this in your own tool you could just use the FH and let the Reader interface suck out the content instead of doing the convert.
			fmt.Printf("\nUploading: %s --->  %v\n\n", hash, item)
			req, err := client.NewRequest("PUT", item.Url, bytes.NewReader(files[hash]))
			// TODO:  headers = { "content-type" => "application/x-binary", "content-md5" => checksum64, "accept" => "application/json" }
			if err != nil {
				fmt.Println(os.Stderr, "Issue this shouldn't happen:", err)
			}

			// post the files
			upload := func() error {
				_, err = client.Do(req, nil)
				return err
			}

			// with exp backoff
			err = upload()
			fmt.Println(os.Stderr, "Issue posting files to the sandbox: ", err)
			// TODO: backoff of 4xx and 5xx doesn't make sense
			// err = backoff.Retry(upload, backoff.NewExponentialBackOff())
			// if err != nil {
			// 	fmt.Println(os.Stderr, "Issue posting files to the sandbox: ", err)
			// }
		}
	}

	// Now lets tell the server we have uploaded all the things.
	sandbox, err := client.Sandboxes.Put(postResp.ID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue commiting sandbox: ", err)
	}
	fmt.Printf("Resulting sandbox %+v\n", sandbox)
}
