package services

import (
    "errors"
    "wallet-microservice/internal/models"
    "wallet-microservice/internal/repositories"
    
    "github.com/google/uuid"
)

type WalletService interface {
    CreateWallet(req models.CreateWalletRequest) (*models.WalletResponse, error)
    GetWallet(id uuid.UUID) (*models.WalletResponse, error)
    GetWalletByUserID(userID string) (*models.WalletResponse, error)
    UpdateWallet(id uuid.UUID, req models.CreateWalletRequest) (*models.WalletResponse, error)
    DeleteWallet(id uuid.UUID) error
    CreditWallet(id uuid.UUID, req models.TransactionRequest) (*models.TransactionResponse, error)
    DebitWallet(id uuid.UUID, req models.TransactionRequest) (*models.TransactionResponse, error)
    GetTransactionHistory(walletID uuid.UUID, page, limit int) ([]models.TransactionResponse, error)
}

type walletService struct {
    walletRepo repositories.WalletRepository
}

func NewWalletService(walletRepo repositories.WalletRepository) WalletService {
    return &walletService{
        walletRepo: walletRepo,
    }
}

func (s *walletService) CreateWallet(req models.CreateWalletRequest) (*models.WalletResponse, error) {
    // Check if wallet already exists for user
    existingWallet, _ := s.walletRepo.GetWalletByUserID(req.UserID)
    if existingWallet != nil {
        return nil, errors.New("wallet already exists for this user")
    }
    
    currency := req.Currency
    if currency == "" {
        currency = "USD"
    }
    
    wallet := &models.Wallet{
        UserID:   req.UserID,
        Balance:  0.0,
        Currency: currency,
    }
    
    err := s.walletRepo.CreateWallet(wallet)
    if err != nil {
        return nil, err
    }
    
    return &models.WalletResponse{
        ID:       wallet.ID,
        UserID:   wallet.UserID,
        Balance:  wallet.Balance,
        Currency: wallet.Currency,
    }, nil
}

func (s *walletService) GetWallet(id uuid.UUID) (*models.WalletResponse, error) {
    wallet, err := s.walletRepo.GetWalletByID(id)
    if err != nil {
        return nil, err
    }
    
    return &models.WalletResponse{
        ID:       wallet.ID,
        UserID:   wallet.UserID,
        Balance:  wallet.Balance,
        Currency: wallet.Currency,
    }, nil
}

func (s *walletService) GetWalletByUserID(userID string) (*models.WalletResponse, error) {
    wallet, err := s.walletRepo.GetWalletByUserID(userID)
    if err != nil {
        return nil, err
    }
    
    return &models.WalletResponse{
        ID:       wallet.ID,
        UserID:   wallet.UserID,
        Balance:  wallet.Balance,
        Currency: wallet.Currency,
    }, nil
}

func (s *walletService) UpdateWallet(id uuid.UUID, req models.CreateWalletRequest) (*models.WalletResponse, error) {
    wallet, err := s.walletRepo.GetWalletByID(id)
    if err != nil {
        return nil, err
    }
    
    wallet.Currency = req.Currency
    if wallet.Currency == "" {
        wallet.Currency = "USD"
    }
    
    err = s.walletRepo.UpdateWallet(wallet)
    if err != nil {
        return nil, err
    }
    
    return &models.WalletResponse{
        ID:       wallet.ID,
        UserID:   wallet.UserID,
        Balance:  wallet.Balance,
        Currency: wallet.Currency,
    }, nil
}

func (s *walletService) DeleteWallet(id uuid.UUID) error {
    return s.walletRepo.DeleteWallet(id)
}

func (s *walletService) CreditWallet(id uuid.UUID, req models.TransactionRequest) (*models.TransactionResponse, error) {
    return s.processTransaction(id, req, models.Credit)
}

func (s *walletService) DebitWallet(id uuid.UUID, req models.TransactionRequest) (*models.TransactionResponse, error) {
    return s.processTransaction(id, req, models.Debit)
}

// func (s *walletService) processTransaction(walletID uuid.UUID, req models.TransactionRequest, transactionType models.TransactionType) (*models.TransactionResponse, error) {
//     // Validate wallet exists
//     wallet, err := s.walletRepo.GetWalletByID(walletID)
//     if err != nil {
//         return nil, err
//     }
    
//     // Update wallet balance
//     err = s.walletRepo.UpdateWalletBalance(walletID, req.Amount, transactionType)
//     if err != nil {
//         return nil, err
//     }
    
//     // Create transaction record
//     transaction := &models.Transaction{
//         WalletID:    wallet.ID,
//         Type:        transactionType,
//         Amount:      req.Amount,
//         Description: req.Description,
//         Reference:   req.Reference,
//     }
    
//     err = s.walletRepo.CreateTransaction(transaction)
//     if err != nil {
//         return nil, err
//     }
    
//     return &models.TransactionResponse{
//         ID:          transaction.ID,
//         WalletID:    transaction.WalletID,
//         Type:        transaction.Type,
//         Amount:      transaction.Amount,
//         Description: transaction.Description,
//         Reference:   transaction.Reference,
//         CreatedAt:   transaction.CreatedAt,
//     }, nil
// }

func (s *walletService) processTransaction(
    walletID uuid.UUID,
    req models.TransactionRequest,
    t models.TransactionType,
) (*models.TransactionResponse, error) {

    // Optional: pre-check that wallet exists to return 404 early;
    // not strictly required, as repo will return not found too.
    if _, err := s.walletRepo.GetWalletByID(walletID); err != nil {
        return nil, err
    }

    txModel := &models.Transaction{
        Description: req.Description,
        Reference:   req.Reference,
    }

    // if err := s.walletRepo.ProcessTransactionWithRollback(walletID, req.Amount, t, txModel); err != nil {
    //     return nil, err
    // }

    return &models.TransactionResponse{
        ID:          txModel.ID,
        WalletID:    txModel.WalletID,
        Type:        txModel.Type,
        Amount:      txModel.Amount,
        Description: txModel.Description,
        Reference:   txModel.Reference,
        CreatedAt:   txModel.CreatedAt,
    }, nil
}


func (s *walletService) GetTransactionHistory(walletID uuid.UUID, page, limit int) ([]models.TransactionResponse, error) {
    if page <= 0 {
        page = 1
    }
    if limit <= 0 || limit > 100 {
        limit = 20
    }
    
    offset := (page - 1) * limit
    
    transactions, err := s.walletRepo.GetTransactionsByWalletID(walletID, limit, offset)
    if err != nil {
        return nil, err
    }
    
    var response []models.TransactionResponse
    for _, tx := range transactions {
        response = append(response, models.TransactionResponse{
            ID:          tx.ID,
            WalletID:    tx.WalletID,
            Type:        tx.Type,
            Amount:      tx.Amount,
            Description: tx.Description,
            Reference:   tx.Reference,
            CreatedAt:   tx.CreatedAt,
        })
    }
    
    return response, nil
}
