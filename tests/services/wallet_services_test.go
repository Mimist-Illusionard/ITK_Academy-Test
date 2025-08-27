package services_test

import (
	"testing"

	enums "itk-academy-test/internal"
	"itk-academy-test/internal/models"
	"itk-academy-test/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockWalletRepo struct {
	createFn        func() (models.Wallet, error)
	updateFn        func(*models.Wallet) (*models.Wallet, error)
	deleteFn        func(id uuid.UUID) error
	getFn           func(id uuid.UUID) (*models.Wallet, error)
	allWalletsFn    func() (*[]models.Wallet, error)
	operateAtomicFn func(id uuid.UUID, fn func(*models.Wallet) error) (*models.Wallet, error)
}

func (m *mockWalletRepo) Create() (models.Wallet, error) {
	return m.createFn()
}
func (m *mockWalletRepo) Update(w *models.Wallet) (*models.Wallet, error) {
	return m.updateFn(w)
}
func (m *mockWalletRepo) Delete(id uuid.UUID) error {
	return m.deleteFn(id)
}
func (m *mockWalletRepo) Get(id uuid.UUID) (*models.Wallet, error) {
	return m.getFn(id)
}
func (m *mockWalletRepo) OperateAtomic(id uuid.UUID, fn func(*models.Wallet) error) (*models.Wallet, error) {
	return m.operateAtomicFn(id, fn)
}
func (m *mockWalletRepo) AllWallets() (*[]models.Wallet, error) {
	return m.allWalletsFn()
}

func TestWalletService_Create(t *testing.T) {
	id := uuid.New()
	mockRepo := &mockWalletRepo{
		createFn: func() (models.Wallet, error) {
			return models.Wallet{ID: id, Balance: 0}, nil
		},
	}
	service := services.New(mockRepo)

	w, err := service.Create()
	assert.NoError(t, err)
	assert.Equal(t, 0, w.Balance)
}

func TestWalletService_Amount(t *testing.T) {
	id := uuid.New()
	mockRepo := &mockWalletRepo{
		getFn: func(got uuid.UUID) (*models.Wallet, error) {
			assert.Equal(t, id, got)
			return &models.Wallet{ID: id, Balance: 100}, nil
		},
	}
	service := services.New(mockRepo)

	amount, err := service.Amount(id)
	assert.NoError(t, err)
	assert.Equal(t, 100, amount)
}

func TestWalletService_Operation_Deposit(t *testing.T) {
	id := uuid.New()

	mockRepo := &mockWalletRepo{
		operateAtomicFn: func(got uuid.UUID, fn func(*models.Wallet) error) (*models.Wallet, error) {
			assert.Equal(t, id, got)
			w := &models.Wallet{ID: id, Balance: 100}
			if err := fn(w); err != nil {
				return nil, err
			}
			return w, nil
		},
	}
	svc := services.New(mockRepo)

	w, err := svc.Operation(id, enums.DEPOSIT, 50)
	assert.NoError(t, err)
	assert.Equal(t, 150, w.Balance)
	assert.Equal(t, id, w.ID)
}

func TestWalletService_Operation_Withdraw_Success(t *testing.T) {
	id := uuid.New()

	mockRepo := &mockWalletRepo{
		operateAtomicFn: func(got uuid.UUID, fn func(*models.Wallet) error) (*models.Wallet, error) {
			assert.Equal(t, id, got)
			w := &models.Wallet{ID: id, Balance: 100}
			if err := fn(w); err != nil {
				return nil, err
			}
			return w, nil
		},
	}
	svc := services.New(mockRepo)

	w, err := svc.Operation(id, enums.WITHDRAW, 50)
	assert.NoError(t, err)
	assert.Equal(t, 50, w.Balance)
	assert.Equal(t, id, w.ID)
}

func TestWalletService_Operation_Withdraw_InsufficientFunds(t *testing.T) {
	id := uuid.New()

	mockRepo := &mockWalletRepo{
		operateAtomicFn: func(got uuid.UUID, fn func(*models.Wallet) error) (*models.Wallet, error) {
			assert.Equal(t, id, got)
			w := &models.Wallet{ID: id, Balance: 30}
			if err := fn(w); err != nil {
				return nil, err
			}
			return w, nil
		},
	}
	svc := services.New(mockRepo)

	w, err := svc.Operation(id, enums.WITHDRAW, 50)
	assert.Nil(t, w)
	assert.EqualError(t, err, "Insufficient funds")
}

func TestWalletService_Operation_InvalidType(t *testing.T) {
	id := uuid.New()

	mockRepo := &mockWalletRepo{
		operateAtomicFn: func(got uuid.UUID, fn func(*models.Wallet) error) (*models.Wallet, error) {
			assert.Equal(t, id, got)
			w := &models.Wallet{ID: id, Balance: 100}
			if err := fn(w); err != nil {
				return nil, err
			}
			return w, nil
		},
	}
	svc := services.New(mockRepo)

	w, err := svc.Operation(id, "HELLO", 10)
	assert.Nil(t, w)
	assert.EqualError(t, err, "Error")
}
