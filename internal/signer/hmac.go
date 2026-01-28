package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"
)

type HmacSigner struct {
	appToken  string
	secretKey string
	now       func() time.Time
}

func New(appToken, secretKey string) *HmacSigner {
	return &HmacSigner{appToken: appToken, secretKey: secretKey, now: time.Now}
}

func (s *HmacSigner) Sign(method, path string, body any) (map[string]string, error) {
	ts := strconv.FormatInt(s.now().Unix(), 10)

	payload := ts + method + path
	if body != nil {
		bs, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		payload += string(bs)
	}
	mac := hmac.New(sha256.New, []byte(s.secretKey))
	mac.Write([]byte(payload))

	return map[string]string{
		"X-App-Token":      s.appToken,
		"X-App-Access-Ts":  ts,
		"X-App-Access-Sig": hex.EncodeToString(mac.Sum(nil)),
	}, nil
}
