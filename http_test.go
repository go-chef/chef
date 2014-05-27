package chef

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var testRequiredHeaders = []string{
	"Accept",
	"X-Ops-Timestamp",
	"X-Ops-Userid",
	"X-Ops-Sign",
	"X-Ops-Content-Hash",
	"X-Ops-Authorization-1",
}

const (
	userid     = "tester"
	requestUrl = "http://localhost:80"
	privateKey = "-----BEGIN RSA PRIVATE KEY-----\nMIICXwIBAAKBgQDAoFRfamHVOqmJkmKyufLqvpPwLGN49a/Ze+RQ3pcwdFdb8sex\nEvr/TYAKEcxs057i8Wuaf5pFt8DFXyYL3iJlFwO30WHmeTv7WsGng2GmlxYKkYMg\nWCt5x3twLahPGzP11KSel7cPy4rzKRvkZP7aLiPIfskJ8kKQ2czCsXYibQIDAQAB\nAoGBALwSzs5qnCMJJ8c+ukcu71LryJ3TeTv9Bjkekgmzi4Kv1Svdm8P0eEUVclJi\nlmobJSMH/LvYotQ3WWxcPlWQCZtgNVWbFfAlsIc39zMOk3lsR9MF5EQIcWZZp3i2\n2h2sR1K/2cx0H+/iU7oeuPtkpGVAihb2iDEd7BK+r7jrfbcBAkEA5kAzqtblhEc4\nUPqrgVOZHiScACT8tHC/r4xUC3VqLmnfcOJKOH1E2XhLjb76IHnLD04yOXvmhS++\n58yzQY0jUQJBANYq+/7PMhJRo8AW/MDI1vOBTToKzcvwcBVZqhY/znqrA3Yg26tu\nM9oqezyc3uIN3HOCQuiZbRRVBZeKmY/r7l0CQQDF1IHQFoXrSpoLkeUL4D0eFgxn\nX2A01O8NsP+BPOf3awYNYpCsyoz+YQphhqY4gwzCYMhsdZVR9/0KAuo9tzuRAkEA\n1JzFoHfHKKJ9osPvVd/MbN8PcLCrD2v5iWiDTyU28VZ20D3cdfqoZUxJHapKJjZG\nhTFrBQjTXhztuTyyKEu7TQJBAIzBLyFcBQdLxor2bH2P2ijU/iAsCxWc5I7VE6zi\n34tYrujX4pAsT+v+06/dMsEtojLIMzffzp11l2zddH66j5g=\n-----END RSA PRIVATE KEY-----"
	publicKey  = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDAoFRfamHVOqmJkmKyufLqvpPw\nLGN49a/Ze+RQ3pcwdFdb8sexEvr/TYAKEcxs057i8Wuaf5pFt8DFXyYL3iJlFwO3\n0WHmeTv7WsGng2GmlxYKkYMgWCt5x3twLahPGzP11KSel7cPy4rzKRvkZP7aLiPI\nfskJ8kKQ2czCsXYibQIDAQAB\n-----END PUBLIC KEY-----"
)

// Gave up trying to implement this myself
// nopCloser came from https://groups.google.com/d/msg/golang-nuts/J-Y4LtdGNSw/wDSYbHWIKj0J
// yay for sharing
// nopCloser creates a io.ReadCloser to satisfy the request.Body input
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func createServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(checkHeader))
}

// publicKeyFromString parses an RSA public key from a string
func publicKeyFromString(key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, fmt.Errorf("block size invalid for '%s'", string(key))
	}
	rsaKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsaKey.(*rsa.PublicKey), nil
}

func makeAuthConfig() (*AuthConfig, error) {
	pk, err := privateKeyFromString([]byte(privateKey))
	if err != nil {
		return nil, err
	}

	ac := &AuthConfig{
		privateKey: pk,
		clientName: userid,
		cryptoHash: crypto.SHA1,
	}
	return ac, nil
}

func TestAuthConfig(t *testing.T) {
	_, err := makeAuthConfig()
	if err != nil {
		t.Error("Failed to create AuthConfig struct from privatekeys and stuff", err)
	}
}

func TestSignRequestNoBody(t *testing.T) {
	ac, err := makeAuthConfig()
	request, err := http.NewRequest("GET", requestUrl, nil)

	for _, hash := range []crypto.Hash{crypto.MD5, crypto.SHA1, crypto.SHA224, crypto.SHA256, crypto.SHA384, crypto.SHA512} {
		err = ac.SignRequest(request, hash)
		if err != nil {
			t.Fatal("failed to generate RequestHeaders")
		}
		count := 0
		for _, requiredHeader := range testRequiredHeaders {
			for header := range request.Header {
				if strings.ToLower(requiredHeader) == strings.ToLower(header) {
					count++
					break
				}
			}
		}
		if count != len(testRequiredHeaders) {
			t.Error("apiRequestHeaders didn't return all of testRequiredHeaders")
		}
	}
}

func TestSignRequestBody(t *testing.T) {
	ac, err := makeAuthConfig()
	if err != nil {
		t.Fatal(err)
	}

	// Gave up trying to implement this myself
	// nopCloser came from https://groups.google.com/d/msg/golang-nuts/J-Y4LtdGNSw/wDSYbHWIKj0J
	// yay for sharing
	requestBody := nopCloser{bytes.NewBufferString("somecoolbodytext")}
	request, err := http.NewRequest("GET", requestUrl, requestBody)

	for _, hash := range []crypto.Hash{crypto.MD5, crypto.SHA1, crypto.SHA224, crypto.SHA256, crypto.SHA384, crypto.SHA512} {
		err = ac.SignRequest(request, hash)
		if err != nil {
			t.Fatal("failed to generate RequestHeaders")
		}
		count := 0
		for _, requiredHeader := range testRequiredHeaders {
			for header := range request.Header {
				if strings.ToLower(requiredHeader) == strings.ToLower(header) {
					count++
					break
				}
			}
		}
		if count != len(testRequiredHeaders) {
			t.Error("apiRequestHeaders didn't return all of testRequiredHeaders")
		}
	}
}

// <3 goiardi
// Test our headers as goiardi would
// https://github.com/ctdk/goiardi/blob/master/authentication/authentication.go
// func checkHeader(user_id string, r *http.Request) string {
func checkHeader(rw http.ResponseWriter, req *http.Request) {
	user_id := req.Header.Get("X-OPS-USERID")
	// Since we don't have a real client or user to check against,
	// we'll just verify that input user = output user
	// user, err := actor.GetReqUser(user_id)
	// if err != nil {
	if user_id != userid {
		fmt.Fprintf(rw, "Failed to authenticate as %s", user_id)
	}

	contentHash := req.Header.Get("X-OPS-CONTENT-HASH")
	if contentHash == "" {
		fmt.Fprintf(rw, "no content hash provided")
	}

	authTimestamp := req.Header.Get("x-ops-timestamp")
	if authTimestamp == "" {
		fmt.Fprintf(rw, "no timestamp header provided")
	}
	// TODO: Will want to implement this later
	//  else {
	// 	// check the time stamp w/ allowed slew
	// 	tok, terr := checkTimeStamp(authTimestamp, config.Config.TimeSlewDur)
	// 	if !tok {
	// 		return terr
	// 	}
	// }

	// Eventually this may be put to some sort of use, but for now just
	// make sure that it's there. Presumably eventually it would be used to
	// use algorithms other than sha1 for hashing the body, or using a
	// different version of the header signing algorithm.
	xopssign := req.Header.Get("x-ops-sign")
	var apiVer string
	var hashChk []string
	if xopssign == "" {
		fmt.Fprintf(rw, "missing X-Ops-Sign header")
	} else {
		re := regexp.MustCompile(`version=(\d+\.\d+)`)
		shaRe := regexp.MustCompile(`algorithm=(\w+)`)
		if verChk := re.FindStringSubmatch(xopssign); verChk != nil {
			apiVer = verChk[1]
			if apiVer != "1.0" && apiVer != "1.1" {
				fmt.Fprintf(rw, "Bad version number '%s' in X-Ops-Header", apiVer)

			}
		} else {
			fmt.Fprintf(rw, "malformed version in X-Ops-Header")
		}

		// if algorithm is missing, it uses sha1. Of course, no other
		// hashing algorithm is supported yet...
		if hashChk = shaRe.FindStringSubmatch(xopssign); hashChk != nil {
			if hashChk[1] != "sha1" {
				fmt.Fprintf(rw, "Unsupported hashing algorithm '%s' specified in X-Ops-Header", hashChk[1])
			}
		}
	}

	var cryptoHash crypto.Hash
	switch hashChk[1] {
	case "md5":
		cryptoHash = crypto.MD5
	case "sha1":
		cryptoHash = crypto.SHA1
	case "sha224":
		cryptoHash = crypto.SHA224
	case "sha256":
		cryptoHash = crypto.SHA256
	case "sha384":
		cryptoHash = crypto.SHA384
	case "sha512":
		cryptoHash = crypto.SHA512
	default:
		fmt.Fprintf(rw, "Invalid crypto hashing algorithm: "+hashChk[1])
		return
	}

	if calcBodyHash(req, cryptoHash) != contentHash {
		fmt.Fprintf(rw, "Content hash did not match hash of request body")

	}

	signedHeaders, sherr := assembleSignedHeader(req)
	if sherr != nil {
		fmt.Fprintf(rw, sherr.Error())
	}

	// signedHeaders are base64 encoded still, we'll need to
	// Decode them
	sig, err := base64.StdEncoding.DecodeString(signedHeaders)
	if err != nil {
		fmt.Fprintf(rw, "Unable to decode signed headers "+err.Error())
	}

	headToCheck := assembleHeaderToCheck(req)
	pubKey, err := publicKeyFromString([]byte(publicKey))

	hash := crypto.SHA1.New()
	hash.Write([]byte(headToCheck))
	hashed := hash.Sum(nil)

	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA1, hashed, sig)
	if err != nil {
		fmt.Fprintf(rw, "Unable to verify signature")
	}
}

func TestRequest(t *testing.T) {
	ac, err := makeAuthConfig()
	server := createServer()
	defer server.Close()

	request, err := http.NewRequest("GET", server.URL, nil)

	err = ac.SignRequest(request, ac.cryptoHash)
	if err != nil {
		t.Fatal("failed to generate RequestHeaders")
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Error("Non 200 return code: " + response.Status)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	bodyStr := buf.String()

	if bodyStr != "" {
		t.Error(bodyStr)
	}

}

func TestRequestToEndpoint(t *testing.T) {
	ac, err := makeAuthConfig()
	server := createServer()
	defer server.Close()

	request, err := http.NewRequest("GET", server.URL+"/clients", nil)

	err = ac.SignRequest(request, ac.cryptoHash)
	if err != nil {
		t.Fatal("failed to generate RequestHeaders")
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Error("Non 200 return code: " + response.Status)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	bodyStr := buf.String()

	if bodyStr != "" {
		t.Error(bodyStr)
	}
}

// More Goiardi <3
func assembleSignedHeader(r *http.Request) (string, error) {
	sHeadStore := make(map[int]string)
	authHeader := regexp.MustCompile(`(?i)^X-Ops-Authorization-(\d+)`)
	for k := range r.Header {
		if c := authHeader.FindStringSubmatch(k); c != nil {
			/* Have to put it into a map first, then sort, in case
			* the headers don't come out in the right order */
			// skipping this error because we shouldn't even be
			// able to get here with something that won't be an
			// integer. Famous last words, I'm sure.
			i, _ := strconv.Atoi(c[1])
			sHeadStore[i] = r.Header.Get(k)
		}
	}
	if len(sHeadStore) == 0 {
		return "", errors.New("No authentication headers found!")
	}

	sH := make([]string, len(sHeadStore))
	sHlimit := len(sH)
	for k, v := range sHeadStore {
		if k > sHlimit {
			return "", errors.New("malformed authentication headers")
		}
		sH[k-1] = v
	}
	signedHeaders := strings.Join(sH, "")

	return signedHeaders, nil
}

func assembleHeaderToCheck(r *http.Request) string {

	// To validate the signature it seems to be very particular
	// Would like to use this loop to generate the content
	// But it causes validation to fail.. so we do it explicitly

	// authHeader := regexp.MustCompile(`(?i)^X-Ops-Authorization-(\d+)`)
	// acceptEncoding := regexp.MustCompile(`(?i)^Accept-Encoding`)
	// userAgent := regexp.MustCompile(`(?i)^User-Agent`)
	//
	// var content string
	// for key, value := range r.Header {
	// 	if !authHeader.MatchString(key) && !acceptEncoding.MatchString(key) && !userAgent.MatchString(key) {
	// 		content += fmt.Sprintf("%s:%s\n", key, value)
	// 	}
	// }
	// return content
	var content string
	content += fmt.Sprintf("%s:%s\n", "Method", r.Header.Get("Method"))
	content += fmt.Sprintf("%s:%s\n", "Hashed Path", r.Header.Get("Hashed Path"))
	content += fmt.Sprintf("%s:%s\n", "Accept", r.Header.Get("Accept"))
	content += fmt.Sprintf("%s:%s\n", "X-Chef-Version", r.Header.Get("X-Chef-Version"))
	content += fmt.Sprintf("%s:%s\n", "X-Ops-Timestamp", r.Header.Get("X-Ops-Timestamp"))
	content += fmt.Sprintf("%s:%s\n", "X-Ops-Userid", r.Header.Get("X-Ops-Userid"))
	content += fmt.Sprintf("%s:%s\n", "X-Ops-Sign", r.Header.Get("X-Ops-Sign"))
	content += fmt.Sprintf("%s:%s\n", "X-Ops-Content-Hash", r.Header.Get("X-Ops-Content-Hash"))
	return content
}

func TestGenerateHash(t *testing.T) {
	// Test all the fun hashing algorithms
	for _, hash := range []crypto.Hash{crypto.MD5, crypto.SHA1, crypto.SHA224, crypto.SHA256, crypto.SHA384, crypto.SHA512} {
		// generateHash should panic if it's unable to generate a hash
		_ = generateHash("hi", hash)
	}
}

func TestGenerateSignatureError(t *testing.T) {
	ac, _ := makeAuthConfig()
	sig, err := generateSignature(ac.privateKey, "hi", crypto.MD4)
	if err == nil {
		t.Error("Successfully generated a signature when we shouldn't have: " + string(sig))
	}
}

func TestRequestError(t *testing.T) {
	ac, err := makeAuthConfig()
	if err != nil {
		t.Fatal(err)
	}

	// Gave up trying to implement this myself
	// nopCloser came from https://groups.google.com/d/msg/golang-nuts/J-Y4LtdGNSw/wDSYbHWIKj0J
	// yay for sharing
	requestBody := nopCloser{bytes.NewBufferString("somecoolbodytext")}
	request, err := http.NewRequest("GET", requestUrl, requestBody)

	err = ac.SignRequest(request, crypto.MD4)
	if err == nil {
		t.Error("Successfully signed a request when we shouldn't have")
	}
}

func TestNewClient(t *testing.T) {
	c, err := NewClient("testclient", privateKey)
	if err != nil {
		t.Error("Couldn't make a valid client...\n", err)
	}
	// simple validation on the created client
	if c.Auth.clientName != "testclient" {
		t.Error("unexpected client name: ", c.Auth.clientName)
	}

	// Bad PEM should be an error
	c, err = NewClient("blah", "not a key")
	if err == nil {
		t.Error("Built a client from a bad key string")
	}
}

func TestMakeRequest(t *testing.T) {
	server := createServer()
	defer server.Close()
	c, _ := NewClient("testclient", privateKey)

	resp, err := c.MakeRequest("GET", server.URL, nil)
	if err != nil {
		t.Error("HRRRM! we tried to make a request but it failed :`( ", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Non 200 return code: ", resp.Status)
	}

	// This should fail
	resp, err = c.MakeRequest("whee", "this will break", nil)
	if err == nil {
		t.Error("This terrible request thing should fail and it didn't")
	}

}
