package models

import (
    "time"
    "github.com/google/uuid"
)

type Wallet struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID    string    `json:"user_id" gorm:"not null;unique;index"`
    Balance   float64   `json:"balance" gorm:"not null;default:0"`
    Currency  string    `json:"currency" gorm:"not null;default:'USD'"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
    ID          uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    WalletID    uuid.UUID       `json:"wallet_id" gorm:"not null;index"`
    Type        TransactionType `json:"type" gorm:"not null"`
    Amount      float64         `json:"amount" gorm:"not null"`
    Description string          `json:"description"`
    Reference   string          `json:"reference"`
    CreatedAt   time.Time       `json:"created_at"`
    Wallet      Wallet          `json:"wallet" gorm:"foreignKey:WalletID"`
}

type TransactionType string

const (
    Credit TransactionType = "CREDIT"
    Debit  TransactionType = "DEBIT"
)

type CreateWalletRequest struct {
    UserID   string  `json:"user_id" binding:"required"`
    Currency string  `json:"currency"`
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
