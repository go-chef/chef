package chef

import (
	"encoding/json"
	"fmt"
	"io"
)

// Reader is a map string interface of arbitrary json
type Reader map[string]interface{}

// Read uses json.Marshal internally to provide a []byte interface to our raw json data
func (b *Reader) Read(p []byte) (size int, err error) {
	fmt.Println(b)
	if buf, err := json.Marshal(&b); err == nil {
		copy(p, buf)
		fmt.Println(fmt.Sprintf("%s", buf))
	}
	return len(p), io.EOF
}
