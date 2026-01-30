package model

type GenerateLinkRequest struct {
	UserID     string // 外部用户唯一标识(建议使用项目名加用户ID)
	LevelName  string // Sumsub 配置的 level 名称
	TTL        int32  // 链接有效期（秒）
	Email      string // 用户邮箱
	Phone      string // 用户手机号
	SuccessURL string // 认证成功跳转地址
	RejectURL  string // 认证拒绝跳转地址
}

// WebhookPayload 是 Sumsub Webhook 回调的核心结构。
//
// 常见 Type（事件类型）示例：
// - applicantCreated：创建 applicant
// - applicantPending：进入审核队列/等待审核
// - applicantPersonalInfoChanged：个人信息变更
type WebhookPayload struct {
	// Type 表示本次回调的事件类型。
	Type string `json:"type"`
	// ApplicantID 是 Sumsub 侧 applicant 唯一标识。
	ApplicantID string `json:"applicantId"`
	// ExternalUserID 是你在创建 applicant / 生成链接时传入的业务侧用户标识。
	ExternalUserID string `json:"externalUserId"`
	// ReviewStatus 是审核流程状态
	ReviewStatus string `json:"reviewStatus"`
}
