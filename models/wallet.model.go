package models

import "time"

type Wallet struct {
	ID         string     `json:"id" gorm:"id"`
	DIDAddress string     `json:"did_address" gorm:"column:did_address"`
	CreatedAt  *time.Time `json:"created_at" gorm:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"deleted_at"`
}

func (r Wallet) TableName() string {
	return "wallets"
}
