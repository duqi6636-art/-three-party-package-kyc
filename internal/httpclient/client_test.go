package httpclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dq/kyc-sdk/kycerrors"
)

func TestGetJSON_UnauthorizedMapped(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("unauthorized"))
	}))
	defer srv.Close()

	cli := New(srv.URL, 1)
	var out any
	err := cli.GetJSON(context.Background(), "/x", nil, &out)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, kycerrors.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got: %v", err)
	}
}

func TestGetJSON_NoContent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	cli := New(srv.URL, 1)
	var out map[string]any
	if err := cli.GetJSON(context.Background(), "/x", nil, &out); err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}

func TestGetJSON_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{invalid"))
	}))
	defer srv.Close()

	cli := New(srv.URL, 1)
	var out map[string]any
	if err := cli.GetJSON(context.Background(), "/x", nil, &out); err == nil {
		t.Fatalf("expected error")
	}
}
