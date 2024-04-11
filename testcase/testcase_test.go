package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = []string{"main","bogusapi"}
	main()
}
