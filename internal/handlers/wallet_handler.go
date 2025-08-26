package handlers

import (
	"itk-academy-test/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	ginEngine.POST("/wallets/", h.Create)
	ginEngine.GET("/wallets/", h.AllWallets)
	ginEngine.GET("/wallets/:id", h.Amount)
	ginEngine.DELETE("/wallets/:id", h.Delete)
}

func (h *WalletHandler) Create(c *gin.Context) {
	wallet, err := h.Service.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't create wallet", "detail": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet created", "wallet": wallet})
}

func (h *WalletHandler) Amount(c *gin.Context) {
	id := c.Param("id")
	walletId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID", "detail": err.Error()})
		return
	}

	amount, err := h.Service.Amount(uint(walletId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "There is error with gettint wallet amount", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"amount": amount})
}

func (h *WalletHandler) Delete(c *gin.Context) {

	id := c.Param("id")
	walletId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID", "detail": err.Error()})
		return
	}

	err = h.Service.Delete(uint(walletId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "There is error with deleting wallet", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet deleted"})
}

func (h *WalletHandler) Operation(c *gin.Context) {
	posts, err := h.Service.AllWallets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get wallets"})
		return
	}

	c.JSON(http.StatusOK, posts)
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
