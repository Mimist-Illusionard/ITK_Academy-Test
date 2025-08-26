package dto

type WalletResponse struct {
	WalletID uint   `json:"walletId"`
	Balance  int    `json:"balance"`
	Message  string `json:"message,omitempty"`
}
