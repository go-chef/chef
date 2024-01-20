package chef

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyMD5Checksum(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "md5-test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tempDir) // clean up

	var (
		// if someone changes the test data,
		// you have to also update the below md5 sum
		testData = []byte("hello\nchef\n")
		filePath = path.Join(tempDir, "dat")
	)
	err = os.WriteFile(filePath, testData, 0644)
	assert.Nil(t, err)
	assert.True(t, verifyMD5Checksum(filePath, "70bda176ac4db06f1f66f96ae0693be1"))
}

func TestFileMD5Checksum(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "md5-test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tempDir) // clean up

	var (
		// if someone changes the test data,
		// you have to also update the below md5 sum
		testData = []byte("hello\nchef\n")
		filePath = path.Join(tempDir, "dat")
	)
	err = os.WriteFile(filePath, testData, 0644)
	assert.Nil(t, err)

	checksum, err := fileMD5Checksum(filePath)
	assert.Nil(t, err)
	assert.Equal(t, "70bda176ac4db06f1f66f96ae0693be1", checksum)
}

func TestMd5Base64Checksum(t *testing.T) {
	result, err := md5Base64Checksum("70bda176ac4db06f1f66f96ae0693be1")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "cL2hdqxNsG8fZvlq4Gk74Q==", result)
}
