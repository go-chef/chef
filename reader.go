package chef

import (
	"bytes"
	"encoding/json"
	"io"
)

// JSONReader handles arbitrary types and synthesizes a streaming encoder for them.
func JSONReader(v interface{}) (io.Reader, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
