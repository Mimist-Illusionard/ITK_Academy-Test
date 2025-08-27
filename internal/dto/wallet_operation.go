package dto

import "github.com/google/uuid"

type WalletOperationRequest struct {
	WalletID      uuid.UUID `json:"valletId" binding:"required"`
	OperationType string    `json:"operationType" binding:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        int       `json:"amount" binding:"required,gt=0"`
}
