package services

import (
	"itk-academy-test/internal/models"
	"itk-academy-test/internal/repository"
)

type WalletService struct {
	repo repository.WalletRepository
}

func New(r repository.WalletRepository) *WalletService {
	return &WalletService{repo: r}
}

func (s *WalletService) Create() (models.Wallet, error) {

	wallet, err := s.repo.Create()

	if err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (s *WalletService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *WalletService) Amount(id uint) (int, error) {
	return s.repo.Amount(id)
}

func (s *WalletService) AllWallets() (*[]models.Wallet, error) {
	return s.repo.AllWallets()
}
