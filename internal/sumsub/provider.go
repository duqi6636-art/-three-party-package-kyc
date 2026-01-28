package sumsub

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dq/kyc-sdk/config"
	"github.com/dq/kyc-sdk/internal/httpclient"
	"github.com/dq/kyc-sdk/internal/signer"
	"github.com/dq/kyc-sdk/kycerrors"
	"github.com/dq/kyc-sdk/model"
)

type Provider struct {
	cfg    *config.Config
	http   *httpclient.Client
	signer *signer.HmacSigner
}

func New(cfg *config.Config) (*Provider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("%w: nil", kycerrors.ErrInvalidConfig)
	}
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, fmt.Errorf("%w: BaseURL required", kycerrors.ErrInvalidConfig)
	}
	if strings.TrimSpace(cfg.AppToken) == "" {
		return nil, fmt.Errorf("%w: AppToken required", kycerrors.ErrInvalidConfig)
	}
	if strings.TrimSpace(cfg.SecretKey) == "" {
		return nil, fmt.Errorf("%w: SecretKey required", kycerrors.ErrInvalidConfig)
	}
	if strings.TrimSpace(cfg.WebhookSecret) == "" {
		return nil, fmt.Errorf("%w: Webhook SecretKey required", kycerrors.ErrInvalidConfig)
	}

	http := httpclient.New(cfg.BaseURL, cfg.TimeoutSec)
	sig := signer.New(cfg.AppToken, cfg.SecretKey)

	return &Provider{
		cfg:    cfg,
		http:   http,
		signer: sig,
	}, nil
}

type applicantDTO struct {
	ID             string `json:"id"`
	ExternalUserID string `json:"externalUserId"`
	Review         struct {
		ReviewStatus string `json:"reviewStatus"`
		ReviewResult struct {
			ReviewAnswer string `json:"reviewAnswer"`
		} `json:"reviewResult"`
	} `json:"review"`
}

func (p *Provider) CreateApplicant(ctx context.Context, userID string) (*model.ApplicantInfo, error) {
	if p == nil {
		return nil, errors.New("nil provider")
	}

	path := "/resources/applicants"
	body := map[string]string{
		"externalUserId": userID,
	}

	headers, err := p.signer.Sign(http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}

	var resp applicantDTO
	if err := p.http.PostJSON(ctx, path, body, headers, &resp); err != nil {
		return nil, err
	}

	return mapApplicant(resp), nil
}

func (p *Provider) GetApplicant(ctx context.Context, applicantID string) (*model.ApplicantInfo, error) {
	if p == nil {
		return nil, errors.New("nil provider")
	}

	path := "/resources/applicants/" + applicantID
	headers, err := p.signer.Sign(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp applicantDTO
	if err := p.http.GetJSON(ctx, path, headers, &resp); err != nil {
		return nil, err
	}

	return mapApplicant(resp), nil
}

type verificationDTO struct {
	URL string `json:"url"`
}

type applicantIdentifiers struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type redirectConfig struct {
	SuccessURL string `json:"successUrl,omitempty"`
	RejectURL  string `json:"rejectUrl,omitempty"`
}

type webSDKLinkRequest struct {
	LevelName            string                `json:"levelName"`
	ExternalUserID       string                `json:"externalUserId"`
	TTLInSecs            int32                 `json:"ttlInSecs"`
	ApplicantIdentifiers *applicantIdentifiers `json:"applicantIdentifiers,omitempty"`
	Redirect             *redirectConfig       `json:"redirect,omitempty"`
}

func (p *Provider) GenerateLink(ctx context.Context, req model.GenerateLinkRequest) (string, error) {
	if p == nil {
		return "", errors.New("nil provider")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return "", errors.New("missing user id")
	}
	if strings.TrimSpace(req.LevelName) == "" {
		return "", errors.New("missing level name")
	}

	ttl := req.TTL
	if ttl <= 0 {
		ttl = 1800
	}

	path := "/resources/sdkIntegrations/levels/-/websdkLink"
	body := webSDKLinkRequest{
		ExternalUserID: req.UserID,
		LevelName:      req.LevelName,
		TTLInSecs:      ttl,
	}
	if strings.TrimSpace(req.Email) != "" || strings.TrimSpace(req.Phone) != "" {
		body.ApplicantIdentifiers = &applicantIdentifiers{
			Email: strings.TrimSpace(req.Email),
			Phone: strings.TrimSpace(req.Phone),
		}
	}
	if strings.TrimSpace(req.SuccessURL) != "" || strings.TrimSpace(req.RejectURL) != "" {
		body.Redirect = &redirectConfig{
			SuccessURL: strings.TrimSpace(req.SuccessURL),
			RejectURL:  strings.TrimSpace(req.RejectURL),
		}
	}

	headers, err := p.signer.Sign(http.MethodPost, path, body)
	if err != nil {
		return "", err
	}

	var resp verificationDTO
	if err := p.http.PostJSON(ctx, path, body, headers, &resp); err != nil {
		return "", err
	}
	if strings.TrimSpace(resp.URL) == "" {
		return "", errors.New("empty link")
	}
	return resp.URL, nil
}

func (p *Provider) VerifyAndParseWebhook(headers http.Header, rawBody []byte) (*model.WebhookPayload, error) {
	if p == nil {
		return nil, errors.New("nil provider")
	}

	sig := strings.TrimSpace(headers.Get("X-Payload-Digest"))
	if sig == "" {
		return nil, errors.New("missing signature")
	}

	verified := verifyWebhookDigest(sig, p.cfg.WebhookSecret, rawBody)
	if !verified {
		return nil, errors.New("invalid signature")
	}

	payload := &model.WebhookPayload{}
	if err := json.Unmarshal(rawBody, payload); err != nil {
		return nil, fmt.Errorf("parse webhook payload: %w", err)
	}

	return payload, nil
}

func verifyWebhookDigest(signature, secretKey string, rawBody []byte) bool {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(rawBody)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	return expectedSignature == signature
}

func mapApplicant(dto applicantDTO) *model.ApplicantInfo {
	return &model.ApplicantInfo{
		UserID:      dto.ExternalUserID,
		ApplicantID: dto.ID,
		Status:      mapStatus(dto.Review.ReviewStatus),
		Result:      mapResult(dto.Review.ReviewResult.ReviewAnswer),
		Provider:    "sumsub",
	}
}

func mapStatus(s string) model.KycStatus {
	switch s {
	case "completed", "reviewed":
		return model.StatusReviewed
	case "pending":
		return model.StatusPending
	default:
		return model.StatusUnknown
	}
}

func mapResult(s string) model.KycResult {
	switch s {
	case "GREEN":
		return model.ResultGreen
	case "RED":
		return model.ResultRed
	case "YELLOW":
		return model.ResultYellow
	default:
		return model.ResultNone
	}
}
