package handlers

import (
	enums "itk-academy-test/internal"
	"itk-academy-test/internal/dto"
	"itk-academy-test/internal/models"
	"itk-academy-test/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler interface {
	Iniitalize(ginEngine *gin.Engine)
}

type WalletHandler struct {
	Service *services.WalletService
}

func New(s *services.WalletService) *WalletHandler {
	return &WalletHandler{Service: s}
}

func (h *WalletHandler) Initialize(ginEngine *gin.Engine) {
	v1 := ginEngine.Group("/api/v1")
	{
		v1.POST("/wallets/", h.Create)
		v1.POST("/wallet/", h.Operation)
		v1.GET("/wallets/", h.AllWallets)
		v1.GET("/wallets/:id", h.Amount)
		v1.DELETE("/wallets/:id", h.Delete)
	}
}

func (h *WalletHandler) Create(c *gin.Context) {
	wallet, err := h.Service.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't create wallet", "detail": err.Error()})
	}

	response := dto.WalletResponse{
		WalletID: wallet.ID,
		Balance:  wallet.Balance,
		Message:  "Wallet created",
	}

	c.JSON(http.StatusOK, response)
}

func (h *WalletHandler) Amount(c *gin.Context) {
	id := c.Param("id")
	walletId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID", "detail": err.Error()})
		return
	}

	amount, err := h.Service.Amount(walletId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "There is error with gettint wallet amount", "detail": err.Error()})
		return
	}

	response := dto.WalletResponse{
		WalletID: walletId,
		Balance:  amount,
		Message:  "",
	}

	c.JSON(http.StatusOK, response)
}

func (h *WalletHandler) Delete(c *gin.Context) {

	id := c.Param("id")
	walletId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID", "detail": err.Error()})
		return
	}

	err = h.Service.Delete(walletId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "There is error with deleting wallet", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet deleted"})
}

func (h *WalletHandler) Operation(c *gin.Context) {
	var request dto.WalletOperationRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet := &models.Wallet{}
	wallet, err = h.Service.Operation(request.WalletID, enums.OperationType(request.OperationType), request.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.WalletResponse{
		WalletID: wallet.ID,
		Balance:  wallet.Balance,
		Message:  "Operation completed successfully",
	}

	c.JSON(http.StatusOK, response)
}

// JUST FOR TESTING
func (h *WalletHandler) AllWallets(c *gin.Context) {
	posts, err := h.Service.AllWallets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get wallets"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
