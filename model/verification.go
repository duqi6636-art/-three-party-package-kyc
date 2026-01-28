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

type WebhookPayload struct {
	Type           string `json:"type"`
	ApplicantID    string `json:"applicantId"`
	ExternalUserID string `json:"externalUserId"`
	InspectionID   string `json:"inspectionId"`
	ReviewStatus   string `json:"reviewStatus"`
	ReviewResult   struct {
		ReviewAnswer string `json:"reviewAnswer"`
	} `json:"reviewResult"`
}
