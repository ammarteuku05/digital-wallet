package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Wallet struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	UserID    string         `json:"user_id" gorm:"not null;index"`
	Balance   float64        `json:"balance" gorm:"type:decimal(15,2);default:0"`
	Currency  string         `json:"currency" gorm:"default:'IDR'"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User *User `json:"-" gorm:"foreignKey:UserID;references:ID"`
}

type WalletTransaction struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	WalletID    string         `json:"wallet_id" gorm:"not null;index"`
	Amount      float64        `json:"amount" gorm:"type:decimal(15,2)"`
	Type        string         `json:"type" gorm:"type:enum('WITHDRAWAL','DEPOSIT');not null"`
	Status      string         `json:"status" gorm:"type:enum('PENDING','COMPLETED','FAILED');default:'PENDING'"`
	Description string         `json:"description" gorm:"null"`
	Metadata    datatypes.JSON `json:"metadata" gorm:"type:json;null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Wallet *Wallet `json:"-" gorm:"foreignKey:WalletID;references:ID"`
}

// TableName specifies the table name for Wallet model
func (Wallet) TableName() string {
	return "wallets"
}

// TableName specifies the table name for WalletTransaction model
func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}
