package chef

import (
	"encoding/json"
	"io"
)

// Reader is a map string interface of arbitrary json
type Reader map[string]interface{}

// Read uses json.Marshal internally to provide a []byte interface to our raw json data
func (b *Reader) Read(p []byte) (size int, err error) {
	if buf, err := json.Marshal(&b); err == nil {
		copy(p, buf)
		return len(p), io.EOF
	}
	return len(p), nil
}
