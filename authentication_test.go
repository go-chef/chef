package chef

import (
	"fmt"
	"testing"
)

const teststr = "Hash this string"
const testsha1 = "hdBcDGYOo5/Q4k2DojVVP1ANs3U="
const testsha256 = "HKxj85/WjYxTHye4B2EPs9UPD8PxhplXZ/tjFucgCj4="

func TestBasicHashStr(t *testing.T) {
	hashOut := HashStr(teststr)
	if hashOut != testsha1 {
		t.Error("Incorrect SHA1 value")
	}
}

func TestBasicHashStr256(t *testing.T) {
	hashOut := HashStr256(teststr)
	fmt.Println(hashOut)
	if hashOut != testsha256 {
		t.Error("Incorrect SHA256 value")
	}
}
