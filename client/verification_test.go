package client

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dq/kyc-sdk/config"
)

func TestClient_GenerateLink(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/resources/sdkIntegrations/levels/-/websdkLink" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		var got map[string]any
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}

		if got["externalUserId"] != "user-1" {
			t.Fatalf("externalUserId mismatch: %v", got["externalUserId"])
		}
		if got["levelName"] != "level-1" {
			t.Fatalf("levelName mismatch: %v", got["levelName"])
		}
		if got["ttlInSecs"] != float64(100) {
			t.Fatalf("ttlInSecs mismatch: %v", got["ttlInSecs"])
		}

		ids, ok := got["applicantIdentifiers"].(map[string]any)
		if !ok {
			t.Fatalf("expected applicantIdentifiers")
		}
		if ids["email"] != "a@b.com" {
			t.Fatalf("email mismatch: %v", ids["email"])
		}

		redirect, ok := got["redirect"].(map[string]any)
		if !ok {
			t.Fatalf("expected redirect")
		}
		if redirect["successUrl"] != "https://ok" {
			t.Fatalf("successUrl mismatch: %v", redirect["successUrl"])
		}

		_ = json.NewEncoder(w).Encode(map[string]string{"url": "https://link"})
	}))
	defer srv.Close()

	cli, err := NewClient(&config.Config{
		BaseURL:   srv.URL,
		AppToken:  "app",
		SecretKey: "secret",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	url, err := cli.GenerateLink(context.Background(), GenerateLinkRequest{
		UserID:     "user-1",
		LevelName:  "level-1",
		TTL:        100,
		Email:      "a@b.com",
		SuccessURL: "https://ok",
	})
	if err != nil {
		t.Fatalf("GenerateLink: %v", err)
	}
	if url != "https://link" {
		t.Fatalf("url mismatch: %s", url)
	}
}

func TestClient_VerifyAndParseWebhook(t *testing.T) {
	raw := []byte(`{"type":"applicantReviewed","applicantId":"a1","externalUserId":"u1","inspectionId":"i1","reviewStatus":"completed","reviewResult":{"reviewAnswer":"GREEN"}}`)
	secret := "secret"

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(raw)
	sig := hex.EncodeToString(mac.Sum(nil))

	cli, err := NewClient(&config.Config{
		BaseURL:       "https://example.com",
		AppToken:      "app",
		SecretKey:     "app-secret",
		WebhookSecret: secret,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	h := http.Header{}
	h.Set("X-Payload-Digest", "sha256="+sig)
	ok, payload, err := cli.VerifyAndParseWebhook(h, raw)
	if err != nil {
		t.Fatalf("VerifyAndParseWebhook: %v", err)
	}
	if !ok {
		t.Fatalf("expected verified")
	}
	if payload == nil || payload.ExternalUserID != "u1" || payload.ReviewResult.ReviewAnswer != "GREEN" {
		t.Fatalf("payload mismatch: %+v", payload)
	}
}

func TestClient_VerifyAndParseWebhook_MissingSignature(t *testing.T) {
	raw := []byte(`{"type":"t","applicantId":"a","externalUserId":"u","inspectionId":"i","reviewStatus":"pending","reviewResult":{"reviewAnswer":"YELLOW"}}`)

	cli, err := NewClient(&config.Config{BaseURL: "https://example.com", AppToken: "app", SecretKey: "secret"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ok, payload, err := cli.VerifyAndParseWebhook(http.Header{}, raw)
	if err != nil {
		t.Fatalf("VerifyAndParseWebhook: %v", err)
	}
	if ok {
		t.Fatalf("expected not verified")
	}
	if payload == nil || payload.ExternalUserID != "u" {
		t.Fatalf("payload mismatch: %+v", payload)
	}
}
