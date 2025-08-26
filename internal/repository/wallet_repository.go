package repository

import (
	"itk-academy-test/internal/models"

	"gorm.io/gorm"
)

type WalletRepository interface {
	Create() (models.Wallet, error)
	Delete(id uint) error

	Deposit(wallet *models.Wallet) error
	Withdraw(wallet *models.Wallet) error

	Amount(id uint) (int, error)
	AllWallets() (*[]models.Wallet, error)
}

type WalletGORMRepository struct {
	DB *gorm.DB
}

func (r *WalletGORMRepository) Create() (models.Wallet, error) {
	wallet := models.Wallet{}
	err := r.DB.Create(&wallet).Error
	return wallet, err
}

func (r *WalletGORMRepository) Delete(id uint) error {
	return r.DB.Delete(models.Wallet{}, id).Error
}

func (r *WalletGORMRepository) Deposit(wallet *models.Wallet) error {
	return nil
}

func (r *WalletGORMRepository) Withdraw(wallet *models.Wallet) error {
	return nil
}

func (r *WalletGORMRepository) Amount(id uint) (int, error) {
	var wallet models.Wallet

	err := r.DB.Where("id = ?", id).Find(&wallet).Error
	if err != nil {
		return 0, err
	}

	return wallet.Balance, nil
}

func (r *WalletGORMRepository) AllWallets() (*[]models.Wallet, error) {
	var wallets []models.Wallet

	result := r.DB.
		Find(&wallets)

	if result.Error != nil {
		return nil, result.Error
	}

	return &wallets, nil
}
