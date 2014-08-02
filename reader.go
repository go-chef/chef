package chef

import (
	"bytes"
	"encoding/json"
)

// JSONReader handles arbitrary types and synthesizes a streaming encoder for them.
func JSONReader(v interface{}) (buf *bytes.Buffer, err error) {
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(v)
	return
}
