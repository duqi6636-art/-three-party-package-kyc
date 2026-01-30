package client

import (
	"context"
	"errors"
	"net/http"

	"github.com/dq/kyc-sdk/model"
)

type GenerateLinkRequest = model.GenerateLinkRequest

type WebhookPayload = model.WebhookPayload

func (c *Client) GenerateLink(ctx context.Context, req GenerateLinkRequest) (string, error) {
	if c == nil || c.provider == nil {
		return "", errors.New("nil client")
	}
	return c.provider.GenerateLink(ctx, req)
}

func (c *Client) VerifyAndParseWebhook(headers http.Header, rawBody []byte) (*WebhookPayload, error) {
	if c == nil || c.provider == nil {
		return nil, errors.New("nil client")
	}
	return c.provider.VerifyAndParseWebhook(headers, rawBody)
}
