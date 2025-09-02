package repositories

import (
    "errors"

    "wallet-microservice/internal/database"
    "wallet-microservice/internal/models"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

type WalletRepository interface {
    CreateWallet(wallet *models.Wallet) error
    GetWalletByID(id uuid.UUID) (*models.Wallet, error)
    GetWalletByUserID(userID string) (*models.Wallet, error)
    UpdateWallet(wallet *models.Wallet) error
    DeleteWallet(id uuid.UUID) error
    CreateTransaction(transaction *models.Transaction) error
    GetTransactionsByWalletID(walletID uuid.UUID, limit, offset int) ([]models.Transaction, error)
    UpdateWalletBalance(walletID uuid.UUID, amount float64, transactionType models.TransactionType) error
}

type walletRepository struct {
    db *gorm.DB
}

func NewWalletRepository() WalletRepository {
    return &walletRepository{
        db: database.DB,
    }
}

func (r *walletRepository) CreateWallet(wallet *models.Wallet) error {
    return r.db.Create(wallet).Error
}

func (r *walletRepository) GetWalletByID(id uuid.UUID) (*models.Wallet, error) {
    var wallet models.Wallet
    err := r.db.First(&wallet, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("wallet not found")
        }
        return nil, err
    }
    return &wallet, nil
}

func (r *walletRepository) GetWalletByUserID(userID string) (*models.Wallet, error) {
    var wallet models.Wallet
    err := r.db.First(&wallet, "user_id = ?", userID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("wallet not found")
        }
        return nil, err
    }
    return &wallet, nil
}

func (r *walletRepository) UpdateWallet(wallet *models.Wallet) error {
    return r.db.Save(wallet).Error
}

func (r *walletRepository) DeleteWallet(id uuid.UUID) error {
    result := r.db.Delete(&models.Wallet{}, "id = ?", id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return errors.New("wallet not found")
    }
    return nil
}

func (r *walletRepository) CreateTransaction(transaction *models.Transaction) error {
    return r.db.Create(transaction).Error
}

func (r *walletRepository) GetTransactionsByWalletID(walletID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
    var transactions []models.Transaction
    err := r.db.Where("wallet_id = ?", walletID).
        Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&transactions).Error
    return transactions, err
}

func (r *walletRepository) UpdateWalletBalance(walletID uuid.UUID, amount float64, transactionType models.TransactionType) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        var wallet models.Wallet
        // Use GORM v2 locking clause for SELECT ... FOR UPDATE
        if err := tx.
            Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&wallet, "id = ?", walletID).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return errors.New("wallet not found")
            }
            return err
        }

        if transactionType == models.Debit {
            if wallet.Balance < amount {
                return errors.New("insufficient balance")
            }
            wallet.Balance -= amount
        } else {
            wallet.Balance += amount
        }

        return tx.Save(&wallet).Error
    })
}
