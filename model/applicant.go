package model

type ApplicantInfo struct {
	UserID      string
	ApplicantID string
	Status      KycStatus
	Result      KycResult
	Provider    string
}

type KycStatus string

const (
	StatusPending  KycStatus = "PENDING"
	StatusReviewed KycStatus = "REVIEWED"
	StatusUnknown  KycStatus = "UNKNOWN"
)

type KycResult string

const (
	ResultGreen  KycResult = "GREEN"
	ResultRed    KycResult = "RED"
	ResultYellow KycResult = "YELLOW"
	ResultNone   KycResult = "NONE"
)
