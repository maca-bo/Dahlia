package verify

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// The version number right now is always v0.
	version = "v0"

	slackRequestTimestampHeader = "X-Slack-Request-TimeStamp"

	slackSignatureHeader = "X-Slack-Signature"
)

func VerifyRequest(r *http.Request) (bool, error) {
	timestamp := r.Header.Get(slackRequestTimestampHeader)
	slackSignature := r.Header.Get(slackSignatureHeader)

	if timestamp == "" {
		return false, fmt.Errorf("slackRequestTimestampHeader is blank.")
	}

	if slackSignature == "" {
		return false, fmt.Errorf("slackSignatureHeader is blank.")
	}

	t, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false, fmt.Errorf("strconv.ParseInt(%s): %v", timestamp, err)
	}

	if ageOK, age := checkTimeStamp(t); !ageOK {
		return false, &OldTimeStumpError{t: age}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, fmt.Errorf("ioutil.ReadAll(%v): %v", r.Body, err)
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	basestring := fmt.Sprintf("%s:%s:%s", version, timestamp, body)
	slackSigningSecret := os.Getenv("SLACK_SIGNING_SECRET")

	signature := getSignature([]byte(basestring), []byte(slackSigningSecret))

	trimmed := strings.TrimPrefix(slackSignature, fmt.Sprintf("%s=", version))
	signatureInHeader, err := hex.DecodeString(trimmed)
	if err != nil {
		return false, fmt.Errorf("hex.DecodeString(%s): %v", trimmed, err)
	}

	return hmac.Equal(signature, signatureInHeader), nil

}

func getSignature(base []byte, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(base)

	return h.Sum(nil)
}

func checkTimeStamp(timestamp int64) (bool, time.Duration) {
	t := time.Since(time.Unix(timestamp, 0))

	return t.Minutes() <= 5, t
}
