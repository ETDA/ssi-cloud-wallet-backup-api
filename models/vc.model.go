package models

import (
	"time"
)

type VC struct {
	ID             string     `json:"id" gorm:"id"`
	CID            string     `json:"cid" gorm:"column:cid"`
	SchemaType     string     `json:"schema_type" gorm:"schema_type"`
	IssuanceDate   *time.Time `json:"issuance_date" gorm:"issuance_date"`
	Issuer         string     `json:"issuer" gorm:"issuer"`
	Holder         string     `json:"holder" gorm:"holder"`
	JWT            string     `json:"jwt" gorm:"jwt"`
}

func (r VC) TableName() string {
	return "vcs"
}
