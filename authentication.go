package chef

import (
	"crypto"
	// Need to include md5 to allow new()
	// on md5
	_ "crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	// Need to include sha1 to allow new()
	// on sha1
	_ "crypto/sha1"
	// Need to include sha256 to allow new()
	// on sha224 and sha256
	_ "crypto/sha256"
	// Need to include sha512 to allow new()
	// On sha384 and sha512
	_ "crypto/sha512"
	"encoding/base64"
	"hash"
)

// Don't export me bro
// generateSignature will generate a signature ( sign ) the given data with the specified crypto.Hash
func generateSignature(priv *rsa.PrivateKey, data string, hash crypto.Hash) (enc []byte, err error) {
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, hash, generateHash(data, hash))
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// generateHash will generate a hash of the specified crypto.Hash of a given string
// and return a []byte
func generateHash(toHash string, hasher crypto.Hash) []byte {
	var hashfunc hash.Hash

	switch hasher {
	case crypto.MD5:
		hashfunc = crypto.MD5.New()
	case crypto.SHA1:
		hashfunc = crypto.SHA1.New()
	case crypto.SHA224:
		hashfunc = crypto.SHA224.New()
	case crypto.SHA256:
		hashfunc = crypto.SHA256.New()
	case crypto.SHA384:
		hashfunc = crypto.SHA384.New()
	case crypto.SHA512:
		hashfunc = crypto.SHA512.New()
	default:
		return nil
	}

	hashfunc.Write([]byte(toHash))
	return hashfunc.Sum(nil)
}

// base64BlockEncode takes a byte slice and breaks it up into a
// slice of base64 encoded strings
func base64BlockEncode(content []byte, limit int) []string {
	resultString := base64.StdEncoding.EncodeToString(content)
	var resultSlice []string

	index := 0

	var maxLengthPerSlice int

	// No limit
	if limit == 0 {
		maxLengthPerSlice = len(resultString)
	} else {
		maxLengthPerSlice = limit
	}

	// Iterate through the encoded string storing
	// a max of <limit> per slice item
	for i := 0; i < len(resultString)/maxLengthPerSlice; i++ {
		resultSlice = append(resultSlice, resultString[index:index+maxLengthPerSlice])
		index += maxLengthPerSlice
	}

	// Add remaining chunk to the end of the slice
	if len(resultString)%maxLengthPerSlice != 0 {
		resultSlice = append(resultSlice, resultString[index:])
	}

	return resultSlice
}
