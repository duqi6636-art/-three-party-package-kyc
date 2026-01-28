package client

import (
	"errors"
	"testing"

	"github.com/dq/kyc-sdk/config"
	"github.com/dq/kyc-sdk/kycerrors"
)

func TestNewClient_InvalidConfig(t *testing.T) {
	if _, err := NewClient(nil); err == nil {
		t.Fatalf("expected error")
	}

	if _, err := NewClient(&config.Config{}); err == nil {
		t.Fatalf("expected error")
	} else if !errors.Is(err, kycerrors.ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got: %v", err)
	}
}
