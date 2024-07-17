//
// Test the go-chef/chef chef server api http code
//

package testapi

import (
	"fmt"
	"github.com/go-chef/chef"
	"log"
	"net/http"
	"os"
)

// http exercise the chef server api and config settings
func Http() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Use config to set a roundtripper function
	cfg := &chef.Config{Timeout: 1, RoundTripper: newTestRt}
	client := Client(cfg)

	// List the current groups
	groupList, err := client.Groups.List()
	if err != nil {
		log.Println("FAILURE listing the existing groups:", err)
	}
	fmt.Fprintln(os.Stdout, "SUCCESS listing the existing groups:", groupList)

	// print roundtripper information
	rt, _ := client.Client.Transport.(*testRt)
	groupList, err = client.Groups.List()
	errcnt := 0
	if err != nil {
		errcnt++
	}
	groupList, err = client.Groups.List()
	if err != nil {
		errcnt++
	}
	if errcnt != rt.err_count {
		log.Printf("FAILURE roundtrip err count is: %+v should be: %+v\n", errcnt, rt.err_count)
	}
	if rt.req_count != 3 {
		log.Printf("FAILURE roundtrip call count is: %+v should be 3\n", rt.req_count)
	}
}

type testRt struct {
	req_count int
	err_count int
	next      http.RoundTripper
}

func newTestRt(next http.RoundTripper) http.RoundTripper { return &testRt{next: next} }

func (this *testRt) RoundTrip(req *http.Request) (*http.Response, error) {
	this.req_count++
	res, err := this.next.RoundTrip(req)
	if err != nil {
		this.err_count++
	}
	return res, err
}
