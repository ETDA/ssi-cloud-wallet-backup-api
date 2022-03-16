package views

import "time"

type VCVerification struct {
	VerificationResult bool        `json:"verification_result"`
	CID                string      `json:"cid"`
	Status             *string     `json:"status"`
	IssuanceDate       *time.Time  `json:"issuance_date"`
	RevokeDate         *time.Time  `json:"revoke_date,omitempty"`
	Type               []string    `json:"type"`
	Issuer             string      `json:"issuer"`
	Holder             string      `json:"holder"`

}
