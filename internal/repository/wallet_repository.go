package repository

import (
	"itk-academy-test/internal/models"

	"gorm.io/gorm"
)

type WalletRepository interface {
	Create() (models.Wallet, error)
	Update(*models.Wallet) (*models.Wallet, error)
	Delete(id uint) error
	Get(id uint) (*models.Wallet, error)

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

func (r *WalletGORMRepository) Update(wallet *models.Wallet) (*models.Wallet, error) {

	err := r.DB.Save(wallet).Error
	if err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (r *WalletGORMRepository) Delete(id uint) error {
	return r.DB.Delete(models.Wallet{}, id).Error
}

func (r *WalletGORMRepository) Get(id uint) (*models.Wallet, error) {
	var wallet models.Wallet

	err := r.DB.First(&wallet, id).Error
	if err != nil {
		return nil, err
	}

	return &wallet, nil
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
