package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid();column:id"`
	UserID    string    `json:"user_id" gorm:"type:varchar(255);not null;uniqueIndex:idx_wallets_user_id;column:user_id"`
	Balance   float64   `json:"balance" gorm:"type:decimal(15,2);not null;default:0.00;column:balance"`
	Currency  string    `json:"currency" gorm:"type:varchar(3);not null;default:'USD';column:currency"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP;column:updated_at"`
}

// TableName specifies the table name for Wallet
func (Wallet) TableName() string {
	return "wallets"
}

// BeforeCreate GORM hook to set ID if not set
func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	if w.CreatedAt.IsZero() {
		w.CreatedAt = time.Now()
	}
	if w.UpdatedAt.IsZero() {
		w.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate GORM hook to set updated_at
func (w *Wallet) BeforeUpdate(tx *gorm.DB) error {
	w.UpdatedAt = time.Now()
	return nil
}

type Transaction struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid();column:id"`
	WalletID    uuid.UUID       `json:"wallet_id" gorm:"type:uuid;not null;index:idx_transactions_wallet_id;column:wallet_id"`
	Type        TransactionType `json:"type" gorm:"type:varchar(10);not null;check:type IN ('CREDIT', 'DEBIT');index:idx_transactions_type;column:type"`
	Amount      float64         `json:"amount" gorm:"type:decimal(15,2);not null;column:amount"`
	Description string          `json:"description" gorm:"type:text;column:description"`
	Reference   string          `json:"reference" gorm:"type:varchar(255);column:reference"`
	CreatedAt   time.Time       `json:"created_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP;index:idx_transactions_created_at;column:created_at"`
	Wallet      Wallet          `json:"wallet" gorm:"foreignKey:WalletID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for Transaction
func (Transaction) TableName() string {
	return "transactions"
}

// BeforeCreate GORM hook to set ID if not set
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	return nil
}

type TransactionType string

const (
	Credit TransactionType = "CREDIT"
	Debit  TransactionType = "DEBIT"
)

type CreateWalletRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Currency string `json:"currency"`
}

type TransactionRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
	Reference   string  `json:"reference"`
}

type WalletResponse struct {
	ID       uuid.UUID `json:"id"`
	UserID   string    `json:"user_id"`
	Balance  float64   `json:"balance"`
	Currency string    `json:"currency"`
}

type TransactionResponse struct {
	ID          uuid.UUID       `json:"id"`
	WalletID    uuid.UUID       `json:"wallet_id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	Description string          `json:"description"`
	Reference   string          `json:"reference"`
	CreatedAt   time.Time       `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
