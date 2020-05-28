package chef

import (
	"testing"
)

const (
	teststr     = "Hash this string"
	testsha1    = "hdBcDGYOo5/Q4k2DojVVP1ANs3U="
	testsha256  = "HKxj85/WjYxTHye4B2EPs9UPD8PxhplXZ/tjFucgCj4="
	privateKeyG = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAx12nDxxOwSPHRSJEDz67a0folBqElzlu2oGMiUTS+dqtj3FU
h5lJc1MjcprRVxcDVwhsSSo9948XEkk39IdblUCLohucqNMzOnIcdZn8zblN7Cnp
W03UwRM0iWX1HuwHnGvm6PKeqKGqplyIXYO0qlDWCzC+VaxFTwOUk31MfOHJQn4y
fTrfuE7h3FTElLBu065SFp3dPICIEmWCl9DadnxbnZ8ASxYQ9xG7hmZduDgjNW5l
3x6/EFkpym+//D6AbWDcVJ1ovCsJL3CfH/NZC3ekeJ/aEeLxP/vaCSH1VYC5VsYK
5Qg7SIa6Nth3+RZz1hYOoBJulEzwljznwoZYRQIDAQABAoIBADPQol+qAsnty5er
PTcdHcbXLJp5feZz1dzSeL0gdxja/erfEJIhg9aGUBs0I55X69VN6h7l7K8PsHZf
MzzJhUL4QJJETOYP5iuVhtIF0I+DTr5Hck/5nYcEv83KAvgjbiL4ZE486IF5awnL
2OE9HtJ5KfhEleNcX7MWgiIHGb8G1jCqu/tH0GI8Z4cNgUrXMbczGwfbN/5Wc0zo
Dtpe0Tec/Fd0DLFwRiAuheakPjlVWb7AGMDX4TyzCXfMpS1ul2jk6nGFk77uQozF
PQUawCRp+mVS4qecgq/WqfTZZbBlW2L18/kpafvsxG8kJ7OREtrb0SloZNFHEc2Q
70GbgKECgYEA6c/eOrI3Uour1gKezEBFmFKFH6YS/NZNpcSG5PcoqF6AVJwXg574
Qy6RatC47e92be2TT1Oyplntj4vkZ3REv81yfz/tuXmtG0AylH7REbxubxAgYmUT
18wUAL4s3TST2AlK4R29KwBadwUAJeOLNW+Rc4xht1galsqQRb4pUzkCgYEA2kj2
vUhKAB7QFCPST45/5q+AATut8WeHnI+t1UaiZoK41Jre8TwlYqUgcJ16Q0H6KIbJ
jlEZAu0IsJxjQxkD4oJgv8n5PFXdc14HcSQ512FmgCGNwtDY/AT7SQP3kOj0Rydg
N02uuRb/55NJ07Bh+yTQNGA+M5SSnUyaRPIAMW0CgYBgVU7grDDzB60C/g1jZk/G
VKmYwposJjfTxsc1a0gLJvSE59MgXc04EOXFNr4a+oC3Bh2dn4SJ2Z9xd1fh8Bur
UwCLwVE3DBTwl2C/ogiN4C83/1L4d2DXlrPfInvloBYR+rIpUlFweDLNuve2pKvk
llU9YGeaXOiHnGoY8iKgsQKBgQDZKMOHtZYhHoZlsul0ylCGAEz5bRT0V8n7QJlw
12+TSjN1F4n6Npr+00Y9ov1SUh38GXQFiLq4RXZitYKu6wEJZCm6Q8YXd1jzgDUp
IyAEHNsrV7Y/fSSRPKd9kVvGp2r2Kr825aqQasg16zsERbKEdrBHmwPmrsVZhi7n
rlXw1QKBgQDBOyUJKQOgDE2u9EHybhCIbfowyIE22qn9a3WjQgfxFJ+aAL9Bg124
fJIEzz43fJ91fe5lTOgyMF5TtU5ClAOPGtlWnXU0e5j3L4LjbcqzEbeyxvP3sn1z
dYkX7NdNQ5E6tcJZuJCGq0HxIAQeKPf3x9DRKzMnLply6BEzyuAC4g==
-----END RSA PRIVATE KEY-----
`
)

var testblock = "Stuff and nonsense to encode  "

func TestGenerateDigestSignature(t *testing.T) {
	pk, _ := PrivateKeyFromString([]byte(privateKeyG))
	_, err := GenerateDigestSignature(pk, teststr)
	if err != nil {
		t.Error("Error generating signature", err)
	}
}

func TestGenerateSignature(t *testing.T) {
	pk, _ := PrivateKeyFromString([]byte(privateKeyG))
	_, err := GenerateSignature(pk, teststr)
	if err != nil {
		t.Error("Error generating signature", err)
	}
}

func TestPrivateEncrypt(t *testing.T) {
	pk, _ := PrivateKeyFromString([]byte(privateKeyG))
	_, err := privateEncrypt(pk, []byte(teststr))
	if err != nil {
		t.Error("Error encrypting", err)
	}
}

func TestBasicHashStr(t *testing.T) {
	hashOut := HashStr(teststr)
	if hashOut != testsha1 {
		t.Error("Incorrect SHA1 value")
	}
}

func TestBasicHashStr256(t *testing.T) {
	hashOut := HashStr256(teststr)
	if hashOut != testsha256 {
		t.Error("Incorrect SHA256 value")
	}
}

func TestBase64BlockEncode(t *testing.T) {
	blockOut := Base64BlockEncode([]byte(testblock+testblock+testblock), 60)
	if len(blockOut) != 2 {
		t.Errorf("Incorrect number of encoded blocks got %+v", len(blockOut))
	}
	if len(blockOut[0]) != 60 {
		t.Errorf("Incorrect length of the encoded block got %+v", len(blockOut))
	}
}
