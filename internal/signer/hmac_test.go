package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"
)

func TestSign_PayloadIncludesBody(t *testing.T) {
	s := &HmacSigner{
		appToken:  "app",
		secretKey: "secret",
		now: func() time.Time {
			return time.Unix(1, 0)
		},
	}

	headersNoBody, err := s.Sign("POST", "/path", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	headersWithBody, err := s.Sign("POST", "/path", map[string]string{"a": "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if headersNoBody["X-App-Access-Sig"] == headersWithBody["X-App-Access-Sig"] {
		t.Fatalf("expected different signatures when body differs")
	}
}

func TestSign_SignatureMatchesExpected(t *testing.T) {
	s := &HmacSigner{
		appToken:  "app",
		secretKey: "secret",
		now: func() time.Time {
			return time.Unix(1, 0)
		},
	}

	headers, err := s.Sign("POST", "/path", map[string]string{"a": "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := "1POST/path" + `{"a":"b"}`
	mac := hmac.New(sha256.New, []byte("secret"))
	mac.Write([]byte(payload))
	want := hex.EncodeToString(mac.Sum(nil))

	if headers["X-App-Access-Sig"] != want {
		t.Fatalf("signature mismatch: want=%s got=%s", want, headers["X-App-Access-Sig"])
	}
}
