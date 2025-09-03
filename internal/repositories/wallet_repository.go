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
	ProcessTransactionWithRollback(walletID uuid.UUID, amount float64, transactionType models.TransactionType, txModel *models.Transaction) error
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

func (r *walletRepository) ProcessTransactionWithRollback(
    walletID uuid.UUID,
    amount float64,
    t models.TransactionType,
    txReq *models.Transaction,
) error {
    // 1) Start a transaction: db.Transaction opens a BEGIN; it will COMMIT on nil or ROLLBACK on error.
    return r.db.Transaction(func(tx *gorm.DB) error {

        // 2) Lock the target wallet row using SELECT ... FOR UPDATE.
        //    This ensures no other transaction can modify this wallet until this one commits/rolls back.
        var wallet models.Wallet
        if err := tx.
            Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&wallet, "id = ?", walletID).Error; err != nil {
            // Returning an error here causes the transaction helper to ROLLBACK.
            return err
        }

        // 3) Business validation and balance math inside the same transaction.
        if t == models.Debit {
            if wallet.Balance < amount {
                // Any error returned triggers a ROLLBACK and discards changes.
                return errors.New("insufficient balance")
            }
            wallet.Balance -= amount
        } else {
            wallet.Balance += amount
        }

        // 4) Persist the new balance; if this fails (e.g., constraint/connection), return error to rollback.
        if err := tx.Save(&wallet).Error; err != nil {
            return err
        }

        // 5) Prepare and insert the transaction record in the same transaction.
        //    If this insert fails, returning error will rollback the earlier balance update as well.
        txReq.WalletID = wallet.ID
        txReq.Type = t
        txReq.Amount = amount

        if err := tx.Create(txReq).Error; err != nil {
            // ROLLBACK: the previous UPDATE will be undone automatically.
            return err
        }

        // 6) Return nil => transaction helper COMMITs both the balance update and the insert atomically.
        return nil
    })
}


func (r *walletRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
