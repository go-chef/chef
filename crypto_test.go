package chef

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

var aesDecryptionTests = []struct {
	cipher   string
	data     string
	key      string
	iv       string
	tag      string
	expected string
}{
	{
		"aes-256-gcm",
		"tg3UG+t43RfVf2hRaPEqU5Vwllu8Uw==\n",
		"3f5eccf62026458b8038a06daecff4a69640ce3c5e4585aace7243d5a7f01b62",
		"MQ2DVUphtxAwdTfo",
		"nsH3unmSGw9hcOG3KpYXpQ==\n",
		`{"json_wrapper":"bar"}`,
	},
}

func TestDecrypt(t *testing.T) {
	for _, dt := range aesDecryptionTests {
		// auth_tag = Base64.decode64(tag)
		authTag, _ := base64.StdEncoding.DecodeString(dt.tag)
		// data = Base64.decode64(data)
		data, _ := base64.StdEncoding.DecodeString(dt.data)
		// iv = Base64.decode64(iv)
		iv, _ := base64.StdEncoding.DecodeString(dt.iv)

		// key = [key].pack('H*')
		key, _ := hex.DecodeString(dt.key)
		// key = OpenSSL::Digest::SHA256.digest(key)
		cryptKey := sha256.Sum256([]byte(key))

		plaintext, err := DecryptValue(cryptKey[:], iv, authTag, data)
		if err != nil {
			t.Fatal(err)
		}

		if string(plaintext) != dt.expected {
			t.Fatalf("Expected: %s, Actual: %s\n", dt.expected, plaintext)
		}
	}
}

func TestEncrypt(t *testing.T) {
	for _, dt := range aesDecryptionTests {
		// auth_tag = Base64.decode64(tag)
		authTag, _ := base64.StdEncoding.DecodeString(dt.tag)
		// data = Base64.decode64(data)
		data, _ := base64.StdEncoding.DecodeString(dt.data)
		// iv = Base64.decode64(iv)
		iv, _ := base64.StdEncoding.DecodeString(dt.iv)

		// key = [key].pack('H*')
		key, _ := hex.DecodeString(dt.key)
		// key = OpenSSL::Digest::SHA256.digest(key)
		cryptKey := sha256.Sum256([]byte(key))

		tag, ciphertext, err := EncryptValue(cryptKey[:], iv, []byte(dt.expected))
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(ciphertext, data) {
			t.Fatalf("ciphertext: [expected: %q, actual: %q]\n", hex.EncodeToString(data), hex.EncodeToString(ciphertext))
		}

		if !bytes.Equal(tag, authTag) {
			t.Fatalf("tag: [expected: %q, actual: %q]\n", hex.EncodeToString(authTag), hex.EncodeToString(tag))
		}
	}
}
