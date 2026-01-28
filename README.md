# kyc-sdk

一个面向多 KYC Provider 的 Go SDK。对业务方只暴露统一的 `client` API，具体厂商实现下沉到 `internal/<provider>`，方便后续替换/新增认证厂商。

## 安装

```bash
go get github.com/dq/kyc-sdk@v0.1.0
```

## 快速开始（默认 Sumsub）

```go
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
		AppToken:      "YOUR_APP_TOKEN",
		SecretKey:     "YOUR_SECRET_KEY",
		WebhookSecret: "YOUR_WEBHOOK_SECRET",
	}

	cli, err := client.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	url, err := cli.GenerateLink(context.Background(), client.GenerateLinkRequest{
		UserID:     "user-123",
		LevelName:  "id-and-liveness",
		TTL:        1800,
		Email:      "a@b.com",
		SuccessURL: "https://example.com/kyc/success",
		RejectURL:  "https://example.com/kyc/reject",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(url)
}
```

更多示例见 [examples/main.go](file:///d:/2026code/kyc-1/examples/main.go)。

## API

- `CreateApplicant(ctx, userID)`：创建 Applicant
- `GetApplicant(ctx, applicantID)`：查询 Applicant
- `GenerateLink(ctx, req)`：生成 WebSDK 链接
- `VerifyAndParseWebhook(headers, rawBody)`：验签并解析 Webhook

请求/回调结构体位于：

- 生成链接请求：`model.GenerateLinkRequest`（对外在 `client.GenerateLinkRequest` 也可直接使用）
- Webhook：`model.WebhookPayload`（对外在 `client.WebhookPayload` 也可直接使用）

## 错误处理

HTTP 4xx/5xx 会返回 `*kycerrors.HTTPError`，可以用 `errors.Is` 做分类判断：

```go
import (
	"errors"
	"github.com/dq/kyc-sdk/kycerrors"
)

if err != nil {
	if errors.Is(err, kycerrors.ErrUnauthorized) {
		// 401/403
	}
	if errors.Is(err, kycerrors.ErrRateLimited) {
		// 429
	}
}
```

## Webhook（验签与解析）

`VerifyAndParseWebhook` 会尝试从 header 中读取 `X-Payload-Digest`，对 `rawBody` 做 HMAC-SHA256 校验并解析 JSON：

```go
verified, payload, err := cli.VerifyAndParseWebhook(r.Header, rawBody)
if err != nil {
	// JSON 解析失败等
}
if !verified {
	// 签名缺失/不匹配：建议记录并按业务策略拒绝或降级处理
}
_ = payload
```

## 多 Provider 扩展

对外 `client` 只依赖一个 `Provider` 接口：

```go
import (
	"context"
	"net/http"
	"github.com/dq/kyc-sdk/model"
)

type Provider interface {
	CreateApplicant(ctx context.Context, userID string) (*model.ApplicantInfo, error)
	GetApplicant(ctx context.Context, applicantID string) (*model.ApplicantInfo, error)
	GenerateLink(ctx context.Context, req model.GenerateLinkRequest) (string, error)
	VerifyAndParseWebhook(headers http.Header, rawBody []byte) (bool, *model.WebhookPayload, error)
}
```

如果要接入新的厂商（例如 Onfido），建议：

- 在 `internal/onfido` 下实现一个 `Provider`
- 业务侧用 `client.New(provider)` 注入，或者在 `client.NewClient` 中切换默认 Provider

## 运行测试

```bash
go test ./...
```

## 版本发布

```bash
git tag v0.1.0
git push origin v0.1.0
```
