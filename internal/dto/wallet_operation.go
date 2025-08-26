package dto

type WalletOperationRequest struct {
	WalletID      uint   `json:"valletId" binding:"required"`
	OperationType string `json:"operationType" binding:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        int    `json:"amount" binding:"required,gt=0"`
}
