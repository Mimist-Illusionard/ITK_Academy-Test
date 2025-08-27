package services

import (
	"errors"
	enums "itk-academy-test/internal"
	"itk-academy-test/internal/models"
	"itk-academy-test/internal/repository"

	"github.com/google/uuid"
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

func (s *WalletService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *WalletService) Amount(id uuid.UUID) (int, error) {

	wallet, err := s.repo.Get(id)
	if err != nil {
		return 0, err
	}

	return wallet.Balance, nil
}

func (s *WalletService) Operation(id uuid.UUID, operationType enums.OperationType, amount int) (*models.Wallet, error) {

	wallet, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	switch operationType {
	case enums.DEPOSIT:
		wallet.Balance += amount
	case enums.WITHDRAW:
		if wallet.Balance < amount {
			return nil, errors.New("Insufficient funds")
		}
		wallet.Balance -= amount
	default:
		return nil, errors.New("Error")
	}

	s.repo.Update(wallet)

	return wallet, nil
}

func (s *WalletService) AllWallets() (*[]models.Wallet, error) {
	return s.repo.AllWallets()
}
