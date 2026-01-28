package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dq/kyc-sdk/kycerrors"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string, timeoutSec int) *Client {
	if timeoutSec == 0 {
		timeoutSec = 10
	}

	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
	}
}

func (c *Client) GetJSON(ctx context.Context, path string, headers map[string]string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return decode(resp, out)
}

func (c *Client) PostJSON(ctx context.Context, path string, body any, headers map[string]string, out any) error {
	bs, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(bs))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return decode(resp, out)
}

func decode(resp *http.Response, out any) error {
	if resp.StatusCode >= 400 {
		body := readBody(resp.Body, 16<<10)
		return &kycerrors.HTTPError{
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}

	if resp.StatusCode == http.StatusNoContent || out == nil {
		return nil
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	return nil
}

func readBody(r io.Reader, limit int64) string {
	bs, err := io.ReadAll(io.LimitReader(r, limit))
	if err != nil || len(bs) == 0 {
		return ""
	}
	return strings.TrimSpace(string(bs))
}
