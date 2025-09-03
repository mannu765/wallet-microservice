//go:build unit
// +build unit

// services/wallet_service_test.go
package services

import (
	"errors"
	"testing"
	"time"

	"wallet-microservice/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for repositories.WalletRepository using testify/mock.
// Adjust method names/signatures to match your actual interface exactly.
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) CreateWallet(w *models.Wallet) error {
	args := m.Called(w)
	// If the mock injects an ID, do it here to mimic DB behavior.
	if args.Error(0) == nil && w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockWalletRepository) GetWalletByID(id uuid.UUID) (*models.Wallet, error) {
	args := m.Called(id)
	if v := args.Get(0); v != nil {
		return v.(*models.Wallet), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) GetWalletByUserID(userID string) (*models.Wallet, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.(*models.Wallet), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) UpdateWallet(w *models.Wallet) error {
	args := m.Called(w)
	return args.Error(0)
}

func (m *MockWalletRepository) DeleteWallet(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockWalletRepository) UpdateWalletBalance(id uuid.UUID, amount float64, t models.TransactionType) error {
	args := m.Called(id, amount, t)
	return args.Error(0)
}

func (m *MockWalletRepository) CreateTransaction(tx *models.Transaction) error {
	args := m.Called(tx)
	if args.Error(0) == nil {
		if tx.ID == uuid.Nil {
			tx.ID = uuid.New()
		}
		if tx.CreatedAt.IsZero() {
			tx.CreatedAt = time.Now()
		}
	}
	return args.Error(0)
}

func (m *MockWalletRepository) GetTransactionsByWalletID(walletID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	args := m.Called(walletID, limit, offset)
	if v := args.Get(0); v != nil {
		return v.([]models.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) ProcessTransactionWithRollback(walletID uuid.UUID, amount float64, t models.TransactionType, txModel *models.Transaction) error {
	args := m.Called(walletID, amount, t, txModel)
	return args.Error(0)
}

func TestCreateWallet(t *testing.T) {
	t.Run("creates with default USD when empty currency", func(t *testing.T) {
		repo := new(MockWalletRepository)
		svc := NewWalletService(repo)

		userID := "user-1"

		// First check for existing wallet returns nil,nil
		repo.On("GetWalletByUserID", userID).Return((*models.Wallet)(nil), nil).Once()
		// CreateWallet should be called with wallet; we allow any pointer and assert fields inside Run.
		repo.On("CreateWallet", mock.AnythingOfType("*models.Wallet")).Run(func(args mock.Arguments) {
			w := args.Get(0).(*models.Wallet)
			assert.Equal(t, userID, w.UserID)
			assert.Equal(t, 0.0, w.Balance)
			assert.Equal(t, "USD", w.Currency)
		}).Return(nil).Once()

		resp, err := svc.CreateWallet(models.CreateWalletRequest{
			UserID:   userID,
			Currency: "",
		})
		assert.NoError(t, err)
		assert.Equal(t, userID, resp.UserID)
		assert.Equal(t, 0.0, resp.Balance)
		assert.Equal(t, "USD", resp.Currency)

		repo.AssertExpectations(t)
	})

	t.Run("returns error if wallet exists", func(t *testing.T) {
		repo := new(MockWalletRepository)
		svc := NewWalletService(repo)
		userID := "user-1"
		repo.On("GetWalletByUserID", userID).Return(&models.Wallet{UserID: userID}, nil).Once()

		resp, err := svc.CreateWallet(models.CreateWalletRequest{UserID: userID})
		assert.Nil(t, resp)
		assert.EqualError(t, err, "wallet already exists for this user")
		repo.AssertExpectations(t)
	})

	t.Run("propagates repository CreateWallet error", func(t *testing.T) {
		repo := new(MockWalletRepository)
		svc := NewWalletService(repo)
		userID := "user-2"
		repo.On("GetWalletByUserID", userID).Return((*models.Wallet)(nil), nil).Once()
		repo.On("CreateWallet", mock.AnythingOfType("*models.Wallet")).Return(errors.New("db error")).Once()

		resp, err := svc.CreateWallet(models.CreateWalletRequest{UserID: userID, Currency: "EUR"})
		assert.Nil(t, resp)
		assert.EqualError(t, err, "db error")
		repo.AssertExpectations(t)
	})
}

func TestGetWallet(t *testing.T) {
	repo := new(MockWalletRepository)
	svc := NewWalletService(repo)

	id := uuid.New()
	w := &models.Wallet{ID: id, UserID: "u", Balance: 10, Currency: "USD"}
	repo.On("GetWalletByID", id).Return(w, nil).Once()

	resp, err := svc.GetWallet(id)
	assert.NoError(t, err)
	assert.Equal(t, id, resp.ID)
	assert.Equal(t, "u", resp.UserID)
	assert.Equal(t, 10.0, resp.Balance)
	assert.Equal(t, "USD", resp.Currency)

	repo.AssertExpectations(t)
}

func TestGetWallet_Error(t *testing.T) {
	repo := new(MockWalletRepository)
	svc := NewWalletService(repo)
	id := uuid.New()
	repo.On("GetWalletByID", id).Return((*models.Wallet)(nil), errors.New("not found")).Once()

	resp, err := svc.GetWallet(id)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "not found")
	repo.AssertExpectations(t)
}

func TestGetWalletByUserID(t *testing.T) {
	repo := new(MockWalletRepository)
	svc := NewWalletService(repo)

	userID := "u1"
	w := &models.Wallet{ID: uuid.New(), UserID: userID, Balance: 5, Currency: "INR"}
	repo.On("GetWalletByUserID", userID).Return(w, nil).Once()

	resp, err := svc.GetWalletByUserID(userID)
	assert.NoError(t, err)
	assert.Equal(t, w.ID, resp.ID)
	assert.Equal(t, userID, resp.UserID)
	assert.Equal(t, 5.0, resp.Balance)
	assert.Equal(t, "INR", resp.Currency)

	repo.AssertExpectations(t)
}

func TestUpdateWallet(t *testing.T) {
	t.Run("sets USD when empty currency", func(t *testing.T) {
		repo := new(MockWalletRepository)
		svc := NewWalletService(repo)
		id := uuid.New()
		w := &models.Wallet{ID: id, UserID: "u", Balance: 1, Currency: "EUR"}

		repo.On("GetWalletByID", id).Return(w, nil).Once()
		repo.On("UpdateWallet", mock.AnythingOfType("*models.Wallet")).Run(func(args mock.Arguments) {
			updated := args.Get(0).(*models.Wallet)
			assert.Equal(t, "USD", updated.Currency)
		}).Return(nil).Once()

		resp, err := svc.UpdateWallet(id, models.CreateWalletRequest{Currency: ""})
		assert.NoError(t, err)
		assert.Equal(t, "USD", resp.Currency)

		repo.AssertExpectations(t)
	})

	t.Run("updates to provided currency", func(t *testing.T) {
		repo := new(MockWalletRepository)
		svc := NewWalletService(repo)
		id := uuid.New()
		w := &models.Wallet{ID: id, UserID: "u", Balance: 1, Currency: "USD"}

		repo.On("GetWalletByID", id).Return(w, nil).Once()
		repo.On("UpdateWallet", mock.AnythingOfType("*models.Wallet")).Return(nil).Once()

		resp, err := svc.UpdateWallet(id, models.CreateWalletRequest{Currency: "GBP"})
		assert.NoError(t, err)
		assert.Equal(t, "GBP", resp.Currency)

		repo.AssertExpectations(t)
	})
}

func TestDeleteWallet(t *testing.T) {
	repo := new(MockWalletRepository)
	svc := NewWalletService(repo)

	id := uuid.New()
	repo.On("DeleteWallet", id).Return(nil).Once()

	err := svc.DeleteWallet(id)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreditDebitWallet(t *testing.T) {
	t.Run("credits wallet and records transaction", func(t *testing.T) {
		repo := new(MockWalletRepository)
		svc := NewWalletService(repo)

		id := uuid.New()
		w := &models.Wallet{ID: id, UserID: "u", Balance: 100, Currency: "USD"}
		req := models.TransactionRequest{Amount: 25, Description: "refund", Reference: "ref-1"}

		repo.On("GetWalletByID", id).Return(w, nil).Once()
		repo.On("ProcessTransactionWithRollback", id, 25.0, models.Credit, mock.AnythingOfType("*models.Transaction")).Run(func(args mock.Arguments) {
			tx := args.Get(3).(*models.Transaction)
			tx.ID = uuid.New()
			tx.WalletID = id
			tx.Type = models.Credit
			tx.Amount = 25.0
			tx.Description = "refund"
			tx.Reference = "ref-1"
			tx.CreatedAt = time.Now()
		}).Return(nil).Once()

		resp, err := svc.CreditWallet(id, req)
		assert.NoError(t, err)
		assert.Equal(t, id, resp.WalletID)
		assert.Equal(t, models.Credit, resp.Type)
		assert.Equal(t, 25.0, resp.Amount)
		assert.NotEqual(t, uuid.Nil, resp.ID)

		repo.AssertExpectations(t)
	})

	t.Run("debits wallet and records transaction", func(t *testing.T) {
		repo := new(MockWalletRepository)
		svc := NewWalletService(repo)

		id := uuid.New()
		w := &models.Wallet{ID: id, UserID: "u", Balance: 100, Currency: "USD"}
		req := models.TransactionRequest{Amount: 40, Description: "purchase", Reference: "order-9"}

		repo.On("GetWalletByID", id).Return(w, nil).Once()
		repo.On("ProcessTransactionWithRollback", id, 40.0, models.Debit, mock.AnythingOfType("*models.Transaction")).Run(func(args mock.Arguments) {
			tx := args.Get(3).(*models.Transaction)
			tx.ID = uuid.New()
			tx.WalletID = id
			tx.Type = models.Debit
			tx.Amount = 40.0
			tx.Description = "purchase"
			tx.Reference = "order-9"
			tx.CreatedAt = time.Now()
		}).Return(nil).Once()

		resp, err := svc.DebitWallet(id, req)
		assert.NoError(t, err)
		assert.Equal(t, models.Debit, resp.Type)
		assert.Equal(t, 40.0, resp.Amount)

		repo.AssertExpectations(t)
	})

	t.Run("propagates errors in sequence", func(t *testing.T) {
		// GetWalletByID error
		{
			repo := new(MockWalletRepository)
			svc := NewWalletService(repo)
			id := uuid.New()
			repo.On("GetWalletByID", id).Return((*models.Wallet)(nil), errors.New("not found")).Once()
			resp, err := svc.CreditWallet(id, models.TransactionRequest{Amount: 1})
			assert.Nil(t, resp)
			assert.EqualError(t, err, "not found")
			repo.AssertExpectations(t)
		}
		// ProcessTransactionWithRollback error (balance update failure)
		{
			repo := new(MockWalletRepository)
			svc := NewWalletService(repo)
			id := uuid.New()
			w := &models.Wallet{ID: id}
			repo.On("GetWalletByID", id).Return(w, nil).Once()
			repo.On("ProcessTransactionWithRollback", id, 1.0, models.Credit, mock.AnythingOfType("*models.Transaction")).Return(errors.New("balance error")).Once()
			resp, err := svc.CreditWallet(id, models.TransactionRequest{Amount: 1})
			assert.Nil(t, resp)
			assert.EqualError(t, err, "balance error")
			repo.AssertExpectations(t)
		}
		// ProcessTransactionWithRollback error (transaction creation failure)
		{
			repo := new(MockWalletRepository)
			svc := NewWalletService(repo)
			id := uuid.New()
			w := &models.Wallet{ID: id}
			repo.On("GetWalletByID", id).Return(w, nil).Once()
			repo.On("ProcessTransactionWithRollback", id, 2.0, models.Credit, mock.AnythingOfType("*models.Transaction")).Return(errors.New("tx error")).Once()
			resp, err := svc.CreditWallet(id, models.TransactionRequest{Amount: 2})
			assert.Nil(t, resp)
			assert.EqualError(t, err, "tx error")
			repo.AssertExpectations(t)
		}
	})
}

func TestGetTransactionHistory(t *testing.T) {
	t.Run("normalizes page and limit; returns transformed responses", func(t *testing.T) {
		repo := new(MockWalletRepository)
		svc := NewWalletService(repo)

		walletID := uuid.New()
		// page <= 0 -> 1, limit <=0 or >100 -> 20
		limit := 20
		offset := 0

		now := time.Now()
		txs := []models.Transaction{
			{ID: uuid.New(), WalletID: walletID, Type: models.Credit, Amount: 10, Description: "a", Reference: "r1", CreatedAt: now},
			{ID: uuid.New(), WalletID: walletID, Type: models.Debit, Amount: 5, Description: "b", Reference: "r2", CreatedAt: now},
		}

		repo.On("GetTransactionsByWalletID", walletID, limit, offset).Return(txs, nil).Once()

		resp, err := svc.GetTransactionHistory(walletID, 0, 0)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		for i := range txs {
			assert.Equal(t, txs[i].ID, resp[i].ID)
			assert.Equal(t, txs[i].Type, resp[i].Type)
		}
		assert.Equal(t, txs[1].Type, resp[1].Type)

		repo.AssertExpectations(t)
	})

	t.Run("propagates repo error", func(t *testing.T) {
		repo := new(MockWalletRepository)
		svc := NewWalletService(repo)

		walletID := uuid.New()
		repo.On("GetTransactionsByWalletID", walletID, 20, 0).Return(nil, errors.New("db")).Once()

		resp, err := svc.GetTransactionHistory(walletID, -1, -10)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "db")

		repo.AssertExpectations(t)
	})
}
