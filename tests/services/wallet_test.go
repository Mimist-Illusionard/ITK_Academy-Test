package services_test

import (
	"testing"

	enums "itk-academy-test/internal"
	"itk-academy-test/internal/models"
	"itk-academy-test/internal/services"

	"github.com/stretchr/testify/assert"
)

type mockWalletRepo struct {
	createFn     func() (models.Wallet, error)
	updateFn     func(*models.Wallet) (*models.Wallet, error)
	deleteFn     func(id uint) error
	getFn        func(id uint) (*models.Wallet, error)
	allWalletsFn func() (*[]models.Wallet, error)
}

func (m *mockWalletRepo) Create() (models.Wallet, error) {
	return m.createFn()
}
func (m *mockWalletRepo) Update(w *models.Wallet) (*models.Wallet, error) {
	return m.updateFn(w)
}
func (m *mockWalletRepo) Delete(id uint) error {
	return m.deleteFn(id)
}
func (m *mockWalletRepo) Get(id uint) (*models.Wallet, error) {
	return m.getFn(id)
}
func (m *mockWalletRepo) AllWallets() (*[]models.Wallet, error) {
	return m.allWalletsFn()
}

func TestWalletService_Create(t *testing.T) {
	mockRepo := &mockWalletRepo{
		createFn: func() (models.Wallet, error) {
			return models.Wallet{ID: 1, Balance: 0}, nil
		},
	}
	service := services.New(mockRepo)

	w, err := service.Create()
	assert.NoError(t, err)
	assert.Equal(t, uint(1), w.ID)
	assert.Equal(t, 0, w.Balance)
}

func TestWalletService_Amount(t *testing.T) {
	mockRepo := &mockWalletRepo{
		getFn: func(id uint) (*models.Wallet, error) {
			return &models.Wallet{ID: id, Balance: 100}, nil
		},
	}
	service := services.New(mockRepo)

	amount, err := service.Amount(1)
	assert.NoError(t, err)
	assert.Equal(t, 100, amount)
}

func TestWalletService_Operation_Deposit(t *testing.T) {
	mockRepo := &mockWalletRepo{
		getFn: func(id uint) (*models.Wallet, error) {
			return &models.Wallet{ID: id, Balance: 100}, nil
		},
		updateFn: func(w *models.Wallet) (*models.Wallet, error) {
			return w, nil
		},
	}
	service := services.New(mockRepo)

	w, err := service.Operation(1, enums.DEPOSIT, 50)
	assert.NoError(t, err)
	assert.Equal(t, 150, w.Balance)
}

func TestWalletService_Operation_Withdraw_Success(t *testing.T) {
	mockRepo := &mockWalletRepo{
		getFn: func(id uint) (*models.Wallet, error) {
			return &models.Wallet{ID: id, Balance: 100}, nil
		},
		updateFn: func(w *models.Wallet) (*models.Wallet, error) {
			return w, nil
		},
	}
	service := services.New(mockRepo)

	w, err := service.Operation(1, enums.WITHDRAW, 50)
	assert.NoError(t, err)
	assert.Equal(t, 50, w.Balance)
}

func TestWalletService_Operation_Withdraw_InsufficientFunds(t *testing.T) {
	mockRepo := &mockWalletRepo{
		getFn: func(id uint) (*models.Wallet, error) {
			return &models.Wallet{ID: id, Balance: 30}, nil
		},
		updateFn: func(w *models.Wallet) (*models.Wallet, error) {
			return w, nil
		},
	}
	service := services.New(mockRepo)

	w, err := service.Operation(1, enums.WITHDRAW, 50)
	assert.Nil(t, w)
	assert.EqualError(t, err, "Insufficient funds")
}

func TestWalletService_Operation_InvalidType(t *testing.T) {
	mockRepo := &mockWalletRepo{
		getFn: func(id uint) (*models.Wallet, error) {
			return &models.Wallet{ID: id, Balance: 100}, nil
		},
		updateFn: func(w *models.Wallet) (*models.Wallet, error) {
			return w, nil
		},
	}
	service := services.New(mockRepo)

	w, err := service.Operation(1, "HELLO", 10)
	assert.Nil(t, w)
	assert.EqualError(t, err, "Error")
}
