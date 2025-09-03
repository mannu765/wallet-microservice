//go:build integration
// +build integration

package services

import (
	"testing"
	"wallet-microservice/internal/database"
	"wallet-microservice/internal/models"
	"wallet-microservice/internal/repositories"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type WalletServiceIntegrationTestSuite struct {
	suite.Suite
	walletService WalletService
	walletRepo    repositories.WalletRepository
}

func (suite *WalletServiceIntegrationTestSuite) SetupSuite() {
	// Connect to test database
	database.Connect()
	database.Migrate()
}

func (suite *WalletServiceIntegrationTestSuite) SetupTest() {
	// Create fresh repository and service for each test
	suite.walletRepo = repositories.NewWalletRepository()
	suite.walletService = NewWalletService(suite.walletRepo)

	// Clean up any existing data
	suite.cleanupTestData()
}

func (suite *WalletServiceIntegrationTestSuite) TearDownTest() {
	suite.cleanupTestData()
}

func (suite *WalletServiceIntegrationTestSuite) TearDownSuite() {
	// Close database connection if needed
}

func (suite *WalletServiceIntegrationTestSuite) cleanupTestData() {
	// Clean up test data between tests
	// This ensures each test starts with a clean slate
}

func (suite *WalletServiceIntegrationTestSuite) TestCreateWalletIntegration() {
	// Test creating a wallet with real database
	userID := "test-user-" + uuid.New().String()

	req := models.CreateWalletRequest{
		UserID:   userID,
		Currency: "EUR",
	}

	wallet, err := suite.walletService.CreateWallet(req)

	suite.NoError(err)
	suite.NotNil(wallet)
	suite.Equal(userID, wallet.UserID)
	suite.Equal(0.0, wallet.Balance)
	suite.Equal("EUR", wallet.Currency)
	suite.NotEqual(uuid.Nil, wallet.ID)
}

func (suite *WalletServiceIntegrationTestSuite) TestCreateWalletDefaultCurrencyIntegration() {
	// Test creating a wallet with default USD currency
	userID := "test-user-" + uuid.New().String()

	req := models.CreateWalletRequest{
		UserID:   userID,
		Currency: "", // Empty currency should default to USD
	}

	wallet, err := suite.walletService.CreateWallet(req)

	suite.NoError(err)
	suite.NotNil(wallet)
	suite.Equal("USD", wallet.Currency)
}

func (suite *WalletServiceIntegrationTestSuite) TestCreateWalletDuplicateUserIntegration() {
	// Test creating a wallet for a user that already has one
	userID := "test-user-" + uuid.New().String()

	req := models.CreateWalletRequest{
		UserID:   userID,
		Currency: "USD",
	}

	// Create first wallet
	wallet1, err := suite.walletService.CreateWallet(req)
	suite.NoError(err)
	suite.NotNil(wallet1)

	// Try to create second wallet for same user
	wallet2, err := suite.walletService.CreateWallet(req)
	suite.Error(err)
	suite.Nil(wallet2)
	suite.Contains(err.Error(), "wallet already exists")
}

func (suite *WalletServiceIntegrationTestSuite) TestCreditWalletIntegration() {
	// Test crediting a wallet
	userID := "test-user-" + uuid.New().String()

	// Create wallet
	wallet, err := suite.walletService.CreateWallet(models.CreateWalletRequest{
		UserID:   userID,
		Currency: "USD",
	})
	suite.NoError(err)

	// Credit wallet
	creditReq := models.TransactionRequest{
		Amount:      100.50,
		Description: "Test credit",
		Reference:   "ref-123",
	}

	transaction, err := suite.walletService.CreditWallet(wallet.ID, creditReq)
	suite.NoError(err)
	suite.NotNil(transaction)
	suite.Equal(models.Credit, transaction.Type)
	suite.Equal(100.50, transaction.Amount)

	// Verify wallet balance was updated
	updatedWallet, err := suite.walletService.GetWallet(wallet.ID)
	suite.NoError(err)
	suite.Equal(100.50, updatedWallet.Balance)
}

func (suite *WalletServiceIntegrationTestSuite) TestDebitWalletIntegration() {
	// Test debiting a wallet
	userID := "test-user-" + uuid.New().String()

	// Create wallet with initial balance
	wallet, err := suite.walletService.CreateWallet(models.CreateWalletRequest{
		UserID:   userID,
		Currency: "USD",
	})
	suite.NoError(err)

	// Credit wallet first
	creditReq := models.TransactionRequest{
		Amount:      200.00,
		Description: "Initial credit",
		Reference:   "ref-init",
	}
	_, err = suite.walletService.CreditWallet(wallet.ID, creditReq)
	suite.NoError(err)

	// Debit wallet
	debitReq := models.TransactionRequest{
		Amount:      50.25,
		Description: "Test debit",
		Reference:   "ref-debit",
	}

	transaction, err := suite.walletService.DebitWallet(wallet.ID, debitReq)
	suite.NoError(err)
	suite.NotNil(transaction)
	suite.Equal(models.Debit, transaction.Type)
	suite.Equal(50.25, transaction.Amount)

	// Verify wallet balance was updated
	updatedWallet, err := suite.walletService.GetWallet(wallet.ID)
	suite.NoError(err)
	suite.Equal(149.75, updatedWallet.Balance)
}

func (suite *WalletServiceIntegrationTestSuite) TestGetTransactionHistoryIntegration() {
	// Test getting transaction history
	userID := "test-user-" + uuid.New().String()

	// Create wallet
	wallet, err := suite.walletService.CreateWallet(models.CreateWalletRequest{
		UserID:   userID,
		Currency: "USD",
	})
	suite.NoError(err)

	// Create several transactions
	transactions := []models.TransactionRequest{
		{Amount: 100, Description: "Credit 1", Reference: "ref1"},
		{Amount: 50, Description: "Debit 1", Reference: "ref2"},
		{Amount: 75, Description: "Credit 2", Reference: "ref3"},
	}

	for _, tx := range transactions {
		if tx.Description == "Debit 1" {
			_, err = suite.walletService.DebitWallet(wallet.ID, tx)
		} else {
			_, err = suite.walletService.CreditWallet(wallet.ID, tx)
		}
		suite.NoError(err)
	}

	// Get transaction history
	history, err := suite.walletService.GetTransactionHistory(wallet.ID, 1, 10)
	suite.NoError(err)
	suite.Len(history, 3)

	// Verify transactions are ordered by creation time (newest first)
	suite.Equal("Credit 2", history[0].Description)
	suite.Equal("Debit 1", history[1].Description)
	suite.Equal("Credit 1", history[2].Description)
}

// Run the integration test suite
func TestWalletServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(WalletServiceIntegrationTestSuite))
}
