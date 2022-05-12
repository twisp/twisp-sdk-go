package token

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func validationErrorTest(t *testing.T, partition string, token string, expectedErr string) {
	t.Helper()

	_, err := NewVerifier("", partition).(tokenVerifier).Verify(token)
	errorContains(t, err, expectedErr)
}

func validationSuccessTest(t *testing.T, partition, token string) {
	t.Helper()
	arn := "arn:aws:iam::123456789012:user/Alice"
	account := "123456789012"
	userID := "Alice"
	_, err := newVerifier(partition, 200, jsonResponse(arn, account, userID), nil).Verify(token)
	if err != nil {
		t.Errorf("received unexpected error: %s", err)
	}
}

func errorContains(t *testing.T, err error, expectedErr string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("err should have contained '%s' was '%s'", expectedErr, err)
	}
}

func assertSTSError(t *testing.T, err error) {
	t.Helper()
	if _, ok := err.(STSError); !ok {
		t.Errorf("Expected err %v to be an STSError but was not", err)
	}
}

var (
	now        = time.Now()
	timeStr    = now.UTC().Format("20060102T150405Z")
	validURL   = fmt.Sprintf("https://sts.amazonaws.com/?action=GetCallerIdentity&X-Amz-Credential=ASIABCDEFGHIJKLMNOPQ%%2F20191216%%2Fus-west-2%%2Fs3%%2Faws4_request&x-amz-signedheaders=x-twisp-aws-id&x-amz-expires=60&x-amz-date=%s", timeStr)
	validToken = toToken(validURL)
)

func toToken(url string) string {
	return v1Prefix + base64.RawURLEncoding.EncodeToString([]byte(url))
}

func newVerifier(partition string, statusCode int, body string, err error) Verifier {
	var rc io.ReadCloser
	if body != "" {
		rc = ioutil.NopCloser(bytes.NewReader([]byte(body)))
	}
	return tokenVerifier{
		client: &http.Client{
			Transport: &roundTripper{
				err: err,
				resp: &http.Response{
					StatusCode: statusCode,
					Body:       rc,
				},
			},
		},
		validSTShostnames: stsHostsForPartition(partition),
	}
}

type roundTripper struct {
	err  error
	resp *http.Response
}

type errorReadCloser struct{}

func (r errorReadCloser) Read(b []byte) (int, error) {
	return 0, errors.New("An Error")
}

func (r errorReadCloser) Close() error {
	return nil
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.resp, rt.err
}

func jsonResponse(arn, account, userid string) string {
	response := getCallerIdentityWrapper{}
	response.GetCallerIdentityResponse.GetCallerIdentityResult.Account = account
	response.GetCallerIdentityResponse.GetCallerIdentityResult.Arn = arn
	response.GetCallerIdentityResponse.GetCallerIdentityResult.UserID = userid
	data, _ := json.Marshal(response)
	return string(data)
}

func TestSTSEndpoints(t *testing.T) {
	cases := []struct {
		partition string
		domain    string
		valid     bool
	}{
		{"aws-cn", "sts.cn-northwest-1.amazonaws.com.cn", true},
		{"aws-cn", "sts.cn-north-1.amazonaws.com.cn", true},
		{"aws-cn", "sts.us-iso-east-1.c2s.ic.gov", false},
		{"aws", "sts.amazonaws.com", true},
		{"aws", "sts-fips.us-west-2.amazonaws.com", true},
		{"aws", "sts-fips.us-east-1.amazonaws.com", true},
		{"aws", "sts.us-east-1.amazonaws.com", true},
		{"aws", "sts.us-east-2.amazonaws.com", true},
		{"aws", "sts.us-west-1.amazonaws.com", true},
		{"aws", "sts.us-west-2.amazonaws.com", true},
		{"aws", "sts.ap-south-1.amazonaws.com", true},
		{"aws", "sts.ap-northeast-1.amazonaws.com", true},
		{"aws", "sts.ap-northeast-2.amazonaws.com", true},
		{"aws", "sts.ap-southeast-1.amazonaws.com", true},
		{"aws", "sts.ap-southeast-2.amazonaws.com", true},
		{"aws", "sts.ca-central-1.amazonaws.com", true},
		{"aws", "sts.eu-central-1.amazonaws.com", true},
		{"aws", "sts.eu-west-1.amazonaws.com", true},
		{"aws", "sts.eu-west-2.amazonaws.com", true},
		{"aws", "sts.eu-west-3.amazonaws.com", true},
		{"aws", "sts.eu-north-1.amazonaws.com", true},
		{"aws", "sts.amazonaws.com.cn", false},
		{"aws", "sts.not-a-region.amazonaws.com", false},
		{"aws-iso", "sts.us-iso-east-1.c2s.ic.gov", true},
		{"aws-iso", "sts.cn-north-1.amazonaws.com.cn", false},
		{"aws-iso-b", "sts.cn-north-1.amazonaws.com.cn", false},
		{"aws-us-gov", "sts.us-gov-east-1.amazonaws.com", true},
		{"aws-us-gov", "sts.amazonaws.com", false},
		{"aws-not-a-partition", "sts.amazonaws.com", false},
	}

	for _, c := range cases {
		verifier := NewVerifier("", c.partition).(tokenVerifier)
		if err := verifier.verifyHost(c.domain); err != nil && c.valid {
			t.Errorf("%s is not valid endpoint for partition %s", c.domain, c.partition)
		}
	}
}

func TestVerifyTokenPreSTSValidations(t *testing.T) {
	b := make([]byte, maxTokenLenBytes+1, maxTokenLenBytes+1)
	s := string(b)
	validationErrorTest(t, "aws", s, "token is too large")
	validationErrorTest(t, "aws", "twisp-aws-v2.asdfasdfa", "token is missing expected \"twisp-aws-v1.\" prefix")
	validationErrorTest(t, "aws", "twisp-aws-v1.decodingerror", "illegal base64 data")

	validationErrorTest(t, "aws", toToken(":ab:cd.af:/asda"), "missing protocol scheme")
	validationErrorTest(t, "aws", toToken("http://"), "unexpected scheme")
	validationErrorTest(t, "aws", toToken("https://google.com"), fmt.Sprintf("unexpected hostname %q in pre-signed URL", "google.com"))
	validationErrorTest(t, "aws-cn", toToken("https://sts.cn-north-1.amazonaws.com.cn/abc"), "unexpected path in pre-signed URL")
	validationErrorTest(t, "aws", toToken("https://sts.amazonaws.com/abc"), "unexpected path in pre-signed URL")
	validationErrorTest(t, "aws", toToken("https://sts.amazonaws.com/?NoInWhiteList=abc"), "non-whitelisted query parameter")
	validationErrorTest(t, "aws", toToken("https://sts.amazonaws.com/?action=get&action=post"), "query parameter with multiple values not supported")
	validationErrorTest(t, "aws", toToken("https://sts.amazonaws.com/?action=NotGetCallerIdenity"), "unexpected action parameter in pre-signed URL")
	validationErrorTest(t, "aws", toToken("https://sts.amazonaws.com/?action=GetCallerIdentity&x-amz-signedheaders=abc%3bx-twisp-aws-i%3bdef"), "client did not sign the x-twisp-aws-id header in the pre-signed URL")
	validationErrorTest(t, "aws", toToken(fmt.Sprintf("https://sts.amazonaws.com/?action=GetCallerIdentity&x-amz-signedheaders=x-twisp-aws-id&x-amz-date=%s&x-amz-expires=9999999", timeStr)), "invalid X-Amz-Expires parameter in pre-signed URL")
	validationErrorTest(t, "aws", toToken("https://sts.amazonaws.com/?action=GetCallerIdentity&x-amz-signedheaders=x-twisp-aws-id&x-amz-date=xxxxxxx&x-amz-expires=60"), "error parsing X-Amz-Date parameter")
	validationErrorTest(t, "aws", toToken("https://sts.amazonaws.com/?action=GetCallerIdentity&x-amz-signedheaders=x-twisp-aws-id&x-amz-date=19900422T010203Z&x-amz-expires=60"), "X-Amz-Date parameter is expired")
	validationErrorTest(t, "aws", toToken(fmt.Sprintf("https://sts.sa-east-1.amazonaws.com/?action=GetCallerIdentity&x-amz-signedheaders=x-twisp-aws-id&x-amz-date=%s&x-amz-expires=60%%gh", timeStr)), "input token was not properly formatted: malformed query parameter")
	validationSuccessTest(t, "aws", toToken(fmt.Sprintf("https://sts.us-east-2.amazonaws.com/?action=GetCallerIdentity&x-amz-signedheaders=x-twisp-aws-id&x-amz-date=%s&x-amz-expires=60", timeStr)))
	validationSuccessTest(t, "aws", toToken(fmt.Sprintf("https://sts.ap-northeast-2.amazonaws.com/?action=GetCallerIdentity&x-amz-signedheaders=x-twisp-aws-id&x-amz-date=%s&x-amz-expires=60", timeStr)))
	validationSuccessTest(t, "aws", toToken(fmt.Sprintf("https://sts.ca-central-1.amazonaws.com/?action=GetCallerIdentity&x-amz-signedheaders=x-twisp-aws-id&x-amz-date=%s&x-amz-expires=60", timeStr)))
	validationSuccessTest(t, "aws", toToken(fmt.Sprintf("https://sts.eu-west-1.amazonaws.com/?action=GetCallerIdentity&x-amz-signedheaders=x-twisp-aws-id&x-amz-date=%s&x-amz-expires=60", timeStr)))
	validationSuccessTest(t, "aws", toToken(fmt.Sprintf("https://sts.sa-east-1.amazonaws.com/?action=GetCallerIdentity&x-amz-signedheaders=x-twisp-aws-id&x-amz-date=%s&x-amz-expires=60", timeStr)))
}

func TestVerifyHTTPError(t *testing.T) {
	_, err := newVerifier("aws", 0, "", errors.New("an error")).Verify(validToken)
	errorContains(t, err, "error during GET: an error")
	assertSTSError(t, err)
}

func TestVerifyHTTP403(t *testing.T) {
	_, err := newVerifier("aws", 403, " ", nil).Verify(validToken)
	errorContains(t, err, "error from AWS (expected 200, got")
	assertSTSError(t, err)
}

func TestVerifyNoRedirectsFollowed(t *testing.T) {
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"UserId":"AROAIIRR6I5NDJBWMIRQQ:admin-session","Account":"111122223333","Arn":"arn:aws:sts::111122223333:assumed-role/Admin/admin-session"}`)
	}))
	defer ts2.Close()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, ts2.URL, http.StatusFound)
	}))
	defer ts.Close()

	tokVerifier := NewVerifier("", "aws").(tokenVerifier)

	resp, err := tokVerifier.client.Get(ts.URL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.Header.Get("Location") != ts2.URL && resp.StatusCode != http.StatusFound {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%#v\n", resp)
		fmt.Println(string(body))
		t.Error("Unexpectedly followed redirect")
	}
}

func TestVerifyBodyReadError(t *testing.T) {
	verifier := tokenVerifier{
		client: &http.Client{
			Transport: &roundTripper{
				err: nil,
				resp: &http.Response{
					StatusCode: 200,
					Body:       errorReadCloser{},
				},
			},
		},
		validSTShostnames: stsHostsForPartition("aws"),
	}
	_, err := verifier.Verify(validToken)
	errorContains(t, err, "error reading HTTP result")
	assertSTSError(t, err)
}

func TestVerifyUnmarshalJSONError(t *testing.T) {
	_, err := newVerifier("aws", 200, "xxxx", nil).Verify(validToken)
	errorContains(t, err, "invalid character")
	assertSTSError(t, err)
}

func TestVerifyInvalidCanonicalARNError(t *testing.T) {
	_, err := newVerifier("aws", 200, jsonResponse("arn", "1000", "userid"), nil).Verify(validToken)
	errorContains(t, err, "arn 'arn' is invalid:")
	assertSTSError(t, err)
}

func TestVerifyInvalidUserIDError(t *testing.T) {
	_, err := newVerifier("aws", 200, jsonResponse("arn:aws:iam::123456789012:user/Alice", "123456789012", "not:vailid:userid"), nil).Verify(validToken)
	errorContains(t, err, "malformed UserID")
	assertSTSError(t, err)
}

func TestVerifyNoSession(t *testing.T) {
	arn := "arn:aws:iam::123456789012:user/Alice"
	account := "123456789012"
	userID := "Alice"
	accessKeyID := "ASIABCDEFGHIJKLMNOPQ"
	identity, err := newVerifier("aws", 200, jsonResponse(arn, account, userID), nil).Verify(validToken)
	if err != nil {
		t.Errorf("expected error to be nil was %q", err)
	}
	if identity.AccessKeyID != accessKeyID {
		t.Errorf("expected AccessKeyID to be %q but was %q", accessKeyID, identity.AccessKeyID)
	}
	if identity.ARN != arn {
		t.Errorf("expected ARN to be %q but was %q", arn, identity.ARN)
	}
	if identity.CanonicalARN != arn {
		t.Errorf("expected CanonicalARN to be %q but was %q", arn, identity.CanonicalARN)
	}
	if identity.UserID != userID {
		t.Errorf("expected Username to be %q but was %q", userID, identity.UserID)
	}
}

func TestVerifySessionName(t *testing.T) {
	arn := "arn:aws:iam::123456789012:user/Alice"
	account := "123456789012"
	userID := "Alice"
	session := "session-name"
	identity, err := newVerifier("aws", 200, jsonResponse(arn, account, userID+":"+session), nil).Verify(validToken)
	if err != nil {
		t.Errorf("expected error to be nil was %q", err)
	}
	if identity.UserID != userID {
		t.Errorf("expected Username to be %q but was %q", userID, identity.UserID)
	}
	if identity.SessionName != session {
		t.Errorf("expected Session to be %q but was %q", session, identity.SessionName)
	}
}

func TestVerifyCanonicalARN(t *testing.T) {
	arn := "arn:aws:sts::123456789012:assumed-role/Alice/extra"
	canonicalARN := "arn:aws:iam::123456789012:role/Alice"
	account := "123456789012"
	userID := "Alice"
	session := "session-name"
	identity, err := newVerifier("aws", 200, jsonResponse(arn, account, userID+":"+session), nil).Verify(validToken)
	if err != nil {
		t.Errorf("expected error to be nil was %q", err)
	}
	if identity.ARN != arn {
		t.Errorf("expected ARN to be %q but was %q", arn, identity.ARN)
	}
	if identity.CanonicalARN != canonicalARN {
		t.Errorf("expected CannonicalARN to be %q but was %q", canonicalARN, identity.CanonicalARN)
	}
}
