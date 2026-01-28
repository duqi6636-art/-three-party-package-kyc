package client

import (
	"context"
	"errors"
	"net/http"

	"github.com/dq/kyc-sdk/config"
	"github.com/dq/kyc-sdk/internal/sumsub"
	"github.com/dq/kyc-sdk/model"
)

type Client struct {
	provider Provider
}

type Provider interface {
	CreateApplicant(ctx context.Context, userID string) (*model.ApplicantInfo, error)
	GetApplicant(ctx context.Context, applicantID string) (*model.ApplicantInfo, error)
	GenerateLink(ctx context.Context, req model.GenerateLinkRequest) (string, error)
	VerifyAndParseWebhook(headers http.Header, rawBody []byte) (*model.WebhookPayload, error)
}

func New(provider Provider) (*Client, error) {
	if provider == nil {
		return nil, errors.New("missing provider")
	}
	return &Client{provider: provider}, nil
}

func NewClient(cfg *config.Config) (*Client, error) {
	p, err := sumsub.New(cfg)
	if err != nil {
		return nil, err
	}
	return New(p)
}
