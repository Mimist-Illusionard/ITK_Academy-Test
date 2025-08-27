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

func (s *WalletService) Operation(id uuid.UUID, op enums.OperationType, amount int) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("Amount must be positive")
	}

	return s.repo.OperateAtomic(id, func(w *models.Wallet) error {
		switch op {
		case enums.DEPOSIT:
			w.Balance += amount
		case enums.WITHDRAW:
			if w.Balance < amount {
				return errors.New("Insufficient funds")
			}
			w.Balance -= amount
		default:
			return errors.New("Error")
		}
		return nil
	})
}

func (s *WalletService) AllWallets() (*[]models.Wallet, error) {
	return s.repo.AllWallets()
}
