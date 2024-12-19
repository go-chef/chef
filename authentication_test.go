package chef

import (
	"crypto/rsa"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

type privateKey struct {
	name string
	key  *rsa.PrivateKey
}

const (
	teststr     = "Hash this string"
	testsha1    = "hdBcDGYOo5/Q4k2DojVVP1ANs3U="
	testsha256  = "HKxj85/WjYxTHye4B2EPs9UPD8PxhplXZ/tjFucgCj4="
	signature10 = "Z6mh9AcCOylVV4mR2MdgQ+FxriOjFipbcowQkDNIOT9o2pVxdzqC4mstZfF430wAEd+HLPXldJafQINmWpHzG2tDd+ms3WQBIq2/kCLmSljP7DHulOsJVeznXtcOwlg6KP4xTK4CYZoh8YuqW3VbdXyGwgx8iOKYac7EBFoKHxu2B4ggDq6NhDvYTA7eNLdXhaotEh6oNkiTByWndUki0boz82R6jMBtMw0ww4ENhZ6Mt3KeC95cFpv5dhmLXWGYVCbCnhWHsUuQlW5Ii1897muSWa+nlQBK6EBiqL8GD/3X8hHta9YLt8zM7W0hqo4NmFKie6mCaaX51Z6ExzUEgw=="
	signature13 = "EuuCVnnbDBah5WiTy9cM1OwhPp+rEQl6L60E9dpYCEcvCMdvs5G6fVuW9wyo7wHlLETMmbltNccdmlByACWKobFi98oH/aYTz6Z7qRxiAJX8Z/b6C6K1c5mxKeJguOt86SGKd6aqQ4O20bXncKJ6u9HxfALphYHIJHegXC414bGbnWbDowI9ZQpNWZ/10bqyoOIiGAOJbzmU2jn6+2eey78sTbKYmddNJzaVzJAq7dtHvGG0yE1h0Y3lCN3mWP9rVu2tVClXKhpOJ1CQGY7t3gVrBtVVvgdQNBW5rhGWDD6BkwRigGiwBczSVk8a1oTiAMRHPXurrEr+c7ArItrUvA=="
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
	privateKey8 = `
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDHXacPHE7BI8dF
IkQPPrtrR+iUGoSXOW7agYyJRNL52q2PcVSHmUlzUyNymtFXFwNXCGxJKj33jxcS
STf0h1uVQIuiG5yo0zM6chx1mfzNuU3sKelbTdTBEzSJZfUe7Aeca+bo8p6ooaqm
XIhdg7SqUNYLML5VrEVPA5STfUx84clCfjJ9Ot+4TuHcVMSUsG7TrlIWnd08gIgS
ZYKX0Np2fFudnwBLFhD3EbuGZl24OCM1bmXfHr8QWSnKb7/8PoBtYNxUnWi8Kwkv
cJ8f81kLd6R4n9oR4vE/+9oJIfVVgLlWxgrlCDtIhro22Hf5FnPWFg6gEm6UTPCW
POfChlhFAgMBAAECggEAM9CiX6oCye3Ll6s9Nx0dxtcsmnl95nPV3NJ4vSB3GNr9
6t8QkiGD1oZQGzQjnlfr1U3qHuXsrw+wdl8zPMmFQvhAkkRM5g/mK5WG0gXQj4NO
vkdyT/mdhwS/zcoC+CNuIvhkTjzogXlrCcvY4T0e0nkp+ESV41xfsxaCIgcZvwbW
MKq7+0fQYjxnhw2BStcxtzMbB9s3/lZzTOgO2l7RN5z8V3QMsXBGIC6F5qQ+OVVZ
vsAYwNfhPLMJd8ylLW6XaOTqcYWTvu5CjMU9BRrAJGn6ZVLip5yCr9ap9NllsGVb
YvXz+Slp++zEbyQns5ES2tvRKWhk0UcRzZDvQZuAoQKBgQDpz946sjdSi6vWAp7M
QEWYUoUfphL81k2lxIbk9yioXoBUnBeDnvhDLpFq0Ljt73Zt7ZNPU7KmWe2Pi+Rn
dES/zXJ/P+25ea0bQDKUftERvG5vECBiZRPXzBQAvizdNJPYCUrhHb0rAFp3BQAl
44s1b5FzjGG3WBqWypBFvilTOQKBgQDaSPa9SEoAHtAUI9JPjn/mr4ABO63xZ4ec
j63VRqJmgrjUmt7xPCVipSBwnXpDQfoohsmOURkC7QiwnGNDGQPigmC/yfk8Vd1z
XgdxJDnXYWaAIY3C0Nj8BPtJA/eQ6PRHJ2A3Ta65Fv/nk0nTsGH7JNA0YD4zlJKd
TJpE8gAxbQKBgGBVTuCsMPMHrQL+DWNmT8ZUqZjCmiwmN9PGxzVrSAsm9ITn0yBd
zTgQ5cU2vhr6gLcGHZ2fhInZn3F3V+HwG6tTAIvBUTcMFPCXYL+iCI3gLzf/Uvh3
YNeWs98ie+WgFhH6silSUXB4Ms2697akq+SWVT1gZ5pc6IecahjyIqCxAoGBANko
w4e1liEehmWy6XTKUIYATPltFPRXyftAmXDXb5NKM3UXifo2mv7TRj2i/VJSHfwZ
dAWIurhFdmK1gq7rAQlkKbpDxhd3WPOANSkjIAQc2ytXtj99JJE8p32RW8anavYq
vzblqpBqyDXrOwRFsoR2sEebA+auxVmGLueuVfDVAoGBAME7JQkpA6AMTa70QfJu
EIht+jDIgTbaqf1rdaNCB/EUn5oAv0GDXbh8kgTPPjd8n3V97mVM6DIwXlO1TkKU
A48a2VaddTR7mPcvguNtyrMRt7LG8/eyfXN1iRfs101DkTq1wlm4kIarQfEgBB4o
9/fH0NErMycumXLoETPK4ALi
-----END PRIVATE KEY-----
`
)

var testblock = "Stuff and nonsense to encode  "

func getPrivateKeys() ([]privateKey, error) {
	pkcs1, err := PrivateKeyFromString([]byte(privateKeyG))
	if err != nil {
		return nil, err
	}
	pkcs8, err := PrivateKeyFromString([]byte(privateKey8))
	if err != nil {
		return nil, err
	}
	return []privateKey{
		{"pkcs#1", pkcs1},
		{"pkcs#8", pkcs8},
	}, nil
}

func TestGenerateDigestSignature(t *testing.T) {
	pks, _ := getPrivateKeys()
	for _, pk := range pks {
		signed, err := GenerateDigestSignature(pk.key, teststr)
		if err != nil {
			t.Error("Error generating signature", pk.name, err)
		}
		assert.Equalf(t, signature13, base64.StdEncoding.EncodeToString(signed),
			"Signature doesn't match. Version=%+v, KeyType=%+v", "1.0", pk.name)
	}
}

func TestGenerateSignature(t *testing.T) {
	pks, _ := getPrivateKeys()
	for _, pk := range pks {
		signed, err := GenerateSignature(pk.key, teststr)
		if err != nil {
			t.Error("Error generating signature", pk.name, err)
		}
		assert.Equalf(t, signature10, base64.StdEncoding.EncodeToString(signed),
			"Signature doesn't match. Version=%+v, KeyType=%+v", "1.3", pk.name)
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
