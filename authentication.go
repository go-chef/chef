package chef

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

// GenerateDigestSignature will generate a signature of the given data protocol 1.3
func GenerateDigestSignature(priv *rsa.PrivateKey, string_to_sign string) (sig []byte, err error) {
	hashed := sha256.Sum256([]byte(string_to_sign))
	sig, err = rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	return sig, nil
}

// GenerateSignature will generate a signature ( sign ) the given data
func GenerateSignature(priv *rsa.PrivateKey, data string) (enc []byte, err error) {
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, 0, []byte(data))
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// HashStr returns the base64 encoded SHA1 sum of the toHash string
func HashStr(toHash string) string {
	h := sha1.New()
	io.WriteString(h, toHash)
	hashed := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return hashed
}

// HashStr256 returns the base64 encoded SHA256 sum of the toHash string
func HashStr256(toHash string) string {
	sum := sha256.Sum256([]byte(toHash))
	sumslice := sum[:]
	hashed := base64.StdEncoding.EncodeToString(sumslice)
	return hashed
}

// Base64BlockEncode takes a byte slice and breaks it up into a
// slice of base64 encoded strings
func Base64BlockEncode(content []byte, limit int) []string {
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
