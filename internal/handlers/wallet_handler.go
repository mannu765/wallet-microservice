package handlers

import (
    "net/http"
    "strconv"
    "wallet-microservice/internal/models"
    "wallet-microservice/internal/services"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type WalletHandler struct {
    walletService services.WalletService
}

func NewWalletHandler(walletService services.WalletService) *WalletHandler {
    return &WalletHandler{
        walletService: walletService,
    }
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
    var req models.CreateWalletRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "validation_error",
            Message: err.Error(),
        })
        return
    }
    
    wallet, err := h.walletService.CreateWallet(req)
    if err != nil {
        c.JSON(http.StatusConflict, models.ErrorResponse{
            Error:   "creation_failed",
            Message: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, wallet)
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "invalid_id",
            Message: "Invalid wallet ID format",
        })
        return
    }
    
    wallet, err := h.walletService.GetWallet(id)
    if err != nil {
        c.JSON(http.StatusNotFound, models.ErrorResponse{
            Error:   "not_found",
            Message: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) GetWalletByUserID(c *gin.Context) {
    userID := c.Param("userId")
    if userID == "" {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "invalid_user_id",
            Message: "User ID is required",
        })
        return
    }
    
    wallet, err := h.walletService.GetWalletByUserID(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, models.ErrorResponse{
            Error:   "not_found",
            Message: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) UpdateWallet(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "invalid_id",
            Message: "Invalid wallet ID format",
        })
        return
    }
    
    var req models.CreateWalletRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "validation_error",
            Message: err.Error(),
        })
        return
    }
    
    wallet, err := h.walletService.UpdateWallet(id, req)
    if err != nil {
        c.JSON(http.StatusNotFound, models.ErrorResponse{
            Error:   "update_failed",
            Message: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) DeleteWallet(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "invalid_id",
            Message: "Invalid wallet ID format",
        })
        return
    }
    
    err = h.walletService.DeleteWallet(id)
    if err != nil {
        c.JSON(http.StatusNotFound, models.ErrorResponse{
            Error:   "deletion_failed",
            Message: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusNoContent, nil)
}

func (h *WalletHandler) CreditWallet(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "invalid_id",
            Message: "Invalid wallet ID format",
        })
        return
    }
    
    var req models.TransactionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "validation_error",
            Message: err.Error(),
        })
        return
    }
    
    transaction, err := h.walletService.CreditWallet(id, req)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "credit_failed",
            Message: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, transaction)
}

func (h *WalletHandler) DebitWallet(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "invalid_id",
            Message: "Invalid wallet ID format",
        })
        return
    }
    
    var req models.TransactionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "validation_error",
            Message: err.Error(),
        })
        return
    }
    
    transaction, err := h.walletService.DebitWallet(id, req)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "debit_failed",
            Message: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, transaction)
}

func (h *WalletHandler) GetTransactionHistory(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "invalid_id",
            Message: "Invalid wallet ID format",
        })
        return
    }
    
    page := 1
    limit := 20
    
    if p := c.Query("page"); p != "" {
        if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
            page = parsed
        }
    }
    
    if l := c.Query("limit"); l != "" {
        if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
            limit = parsed
        }
    }
    
    transactions, err := h.walletService.GetTransactionHistory(id, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Error:   "fetch_failed",
            Message: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "transactions": transactions,
        "page":         page,
        "limit":        limit,
    })
}

func (h *WalletHandler) RegisterRoutes(router *gin.Engine) {
    api := router.Group("/api/v1")
    {
        wallets := api.Group("/wallets")
        {
            wallets.POST("", h.CreateWallet)
            wallets.GET("/:id", h.GetWallet)
            wallets.PUT("/:id", h.UpdateWallet)
            wallets.DELETE("/:id", h.DeleteWallet)
            wallets.POST("/:id/credit", h.CreditWallet)
            wallets.POST("/:id/debit", h.DebitWallet)
            wallets.GET("/:id/transactions", h.GetTransactionHistory)
        }
        
        users := api.Group("/users")
        {
            users.GET("/:userId/wallet", h.GetWalletByUserID)
        }
    }
}
