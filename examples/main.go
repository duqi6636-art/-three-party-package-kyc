package main

import (
	"context"
	"fmt"

	"github.com/dq/kyc-sdk/client"
	"github.com/dq/kyc-sdk/config"
)

func main() {
	cfg := &config.Config{
		BaseURL:       "https://api.sumsub.com",
		AppToken:      "sbx:9MjMj2UIBmYA7AZzgn70f4VB.RzzRbgAW8Ckex24KUuCduSmJrlhxmDIw",
		SecretKey:     "j7ztQmVlNJRzN0vVomTeKAwDbznk9fAE",
		WebhookSecret: "",
	}

	cli, err := client.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	url, err := cli.GenerateLink(context.Background(), client.GenerateLinkRequest{
		UserID:    "user-123",
		LevelName: "id-and-liveness",
		TTL:       1800,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("VerificationLink (custom): %s\n", url)
}
