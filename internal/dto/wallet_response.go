package dto

import "github.com/google/uuid"

type WalletResponse struct {
	WalletID uuid.UUID `json:"walletId"`
	Balance  int       `json:"balance"`
	Message  string    `json:"message,omitempty"`
}
