package chef

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Generates an MD5 checksum from a file at filePath
func fileMD5Checksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	md5String := fmt.Sprintf("%x", hash.Sum(nil))

	return md5String, nil
}

// Verifies the MD5 checksum for a file at filePath
func verifyMD5Checksum(filePath, checksum string) bool {
	md5String, err := fileMD5Checksum(filePath)
	if err != nil {
		return false
	}
	return md5String == checksum
}

// Converts the MD5 checksum to base64 for compatibility with the sandbox upload API
func md5Base64Checksum(md5Checksum string) (string, error) {
	bytes, err := hex.DecodeString(md5Checksum)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
