package client

import (
	"context"
	"errors"

	"github.com/dq/kyc-sdk/model"
)

func (c *Client) CreateApplicant(ctx context.Context, userID string) (*model.ApplicantInfo, error) {
	if c == nil || c.provider == nil {
		return nil, errors.New("nil client")
	}
	return c.provider.CreateApplicant(ctx, userID)
}

func (c *Client) GetApplicant(ctx context.Context, applicantID string) (*model.ApplicantInfo, error) {
	if c == nil || c.provider == nil {
		return nil, errors.New("nil client")
	}
	return c.provider.GetApplicant(ctx, applicantID)
}
