package chef

// https://github.com/golang/go/issues/23514

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"io"
)

// NewCrypter en/decrypts data parsing using aes
func newCrypter(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

func iv() ([]byte, error) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce[:], nil
}

func generateSecret() *[]byte {
	secret := make([]byte, 32)
	rand.Read(secret)
	return &secret
}

// EncryptValue data
func EncryptValue(key, iv, plaintext []byte) (tag []byte, ciphertext []byte, encryptError error) {
	cipher, err := newCrypter(key)
	if err != nil {
		return nil, nil, err
	}

	// https://sourcegraph.com/github.com/golang/go@52cc9e3762171bd45368b3280554bf12a63f23b2/-/blob/src/crypto/aes/aes_gcm.go#L120-131
	out := cipher.Seal(make([]byte, 0), iv, plaintext, nil)

	return out[len(out)-cipher.Overhead():], out[:len(out)-cipher.Overhead()], nil
}

// DecryptValue data
func DecryptValue(key, iv, tag, ciphertext []byte) ([]byte, error) {
	// https://sourcegraph.com/github.com/golang/go@52cc9e3762171bd45368b3280554bf12a63f23b2/-/blob/src/crypto/aes/aes_gcm.go#L155
	cipher, err := newCrypter(key)
	if err != nil {
		return nil, err
	}
	macData := append(ciphertext, tag...)

	return cipher.Open(nil, iv, macData, nil)
}

// EncodeSharedSecret encrypts secret with the the specified and key returns a base64 encoded string
func EncodeSharedSecret(key *rsa.PrivateKey, secret []byte) (string, error) {
	encryptedSecret, err := rsa.EncryptPKCS1v15(rand.Reader, &key.PublicKey, secret)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedSecret), nil
}

// DecodeSharedSecret returns the plaintext shared secret encrypted by the key
func DecodeSharedSecret(key *rsa.PrivateKey, encodedSecret string) ([]byte, error) {
	// Decode encrypted shared secret
	encryptedSecret, err := base64.StdEncoding.DecodeString(encodedSecret)
	if err != nil {
		return nil, err
	}

	// Decrypt shared secret
	return rsa.DecryptPKCS1v15(rand.Reader, key, encryptedSecret)
}
