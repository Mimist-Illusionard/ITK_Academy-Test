package repository

import (
	"itk-academy-test/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	Create() (models.Wallet, error)
	Update(*models.Wallet) (*models.Wallet, error)
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) (*models.Wallet, error)

	AllWallets() (*[]models.Wallet, error)
	OperateAtomic(id uuid.UUID, fn func(w *models.Wallet) error) (*models.Wallet, error)
}

type WalletGORMRepository struct {
	DB *gorm.DB
}

func (r *WalletGORMRepository) Create() (models.Wallet, error) {
	wallet := models.Wallet{ID: uuid.New()}
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

func (r *WalletGORMRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(models.Wallet{}, id).Error
}

func (r *WalletGORMRepository) Get(id uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet

	err := r.DB.First(&wallet, id).Error
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletGORMRepository) OperateAtomic(id uuid.UUID, fn func(w *models.Wallet) error) (*models.Wallet, error) {
	var result *models.Wallet
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var w models.Wallet

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&w, "id = ?", id).Error; err != nil {
			return err
		}

		if err := fn(&w); err != nil {
			return err
		}

		if err := tx.Save(&w).Error; err != nil {
			return err
		}
		result = &w
		return nil
	})
	return result, err
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
