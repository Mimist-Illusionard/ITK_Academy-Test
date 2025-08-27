package repositories_test

import (
	"testing"

	"itk-academy-test/internal/models"
	"itk-academy-test/internal/repository"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresDB *embeddedpostgres.EmbeddedPostgres

func setupTestDB(t *testing.T) *gorm.DB {
	if postgresDB == nil {
		postgresDB = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
			Port(5433).
			Username("postgres").
			Password("postgres").
			Database("test_db"),
		)
		err := postgresDB.Start()
		if err != nil {
			t.Fatalf("failed to start embedded postgres: %v", err)
		}
		t.Cleanup(func() {
			postgresDB.Stop()
		})
	}

	dsn := "host=localhost port=5433 user=postgres password=postgres dbname=test_db sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.Wallet{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}
func TestWalletRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := &repository.WalletGORMRepository{DB: db}

	wallet, err := repo.Create()
	assert.NoError(t, err)
	assert.NotZero(t, wallet.ID)

	got, err := repo.Get(wallet.ID)
	assert.NoError(t, err)
	assert.Equal(t, wallet.ID, got.ID)

	got.Balance = 100
	updated, err := repo.Update(got)
	assert.NoError(t, err)
	assert.Equal(t, int(100), updated.Balance)

	all, err := repo.AllWallets()
	assert.NoError(t, err)
	assert.Len(t, *all, 1)

	err = repo.Delete(wallet.ID)
	assert.NoError(t, err)

	_, err = repo.Get(wallet.ID)
	assert.Error(t, err)
}
