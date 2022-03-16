package views

import "time"

type VCStatus struct {
	CID         string     `json:"cid"`
	DIDAddress  string     `json:"did_address"`
	Status      *string    `json:"status"`
	ActivatedAt *time.Time `json:"activated_at"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
	ExpiredAt   *time.Time `json:"expired_at"`
}
