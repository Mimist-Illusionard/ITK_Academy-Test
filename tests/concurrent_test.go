package tests

import (
	enums "itk-academy-test/internal"
	"itk-academy-test/internal/models"
	"itk-academy-test/internal/repository"
	"itk-academy-test/internal/services"
	"os"
	"sync"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var epg *embeddedpostgres.EmbeddedPostgres

func TestMain(m *testing.M) {
	epg = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(5435).
		Username("postgres").
		Password("postgres").
		Database("test_db"),
	)
	if err := epg.Start(); err != nil {
		panic(err)
	}
	code := m.Run()
	_ = epg.Stop()
	os.Exit(code)
}

func newDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "host=localhost port=5435 user=postgres password=postgres dbname=test_db sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect DB: %v", err)
	}

	err = db.Migrator().DropTable(&models.Wallet{})
	if err != nil {
		t.Fatalf("drop table: %v", err)
	}
	if err := db.AutoMigrate(&models.Wallet{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(200)
		sqlDB.SetMaxIdleConns(200)
	}

	return db
}

func TestConcurrent_Deposit_1000(t *testing.T) {
	db := newDB(t)
	repo := &repository.WalletGORMRepository{DB: db}
	svc := services.New(repo)

	w, err := repo.Create()
	require.NoError(t, err)

	const workers = 1000
	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.Operation(w.ID, enums.DEPOSIT, 1)
			errCh <- err
		}()
	}
	wg.Wait()
	close(errCh)

	for e := range errCh {
		require.NoError(t, e)
	}

	got, err := repo.Get(w.ID)
	require.NoError(t, err)
	assert.Equal(t, 1000, got.Balance)
}

func TestConcurrent_Withdraw_Exact(t *testing.T) {
	db := newDB(t)
	repo := &repository.WalletGORMRepository{DB: db}
	svc := services.New(repo)

	w, err := repo.Create()
	require.NoError(t, err)

	_, err = svc.Operation(w.ID, enums.DEPOSIT, 1000)
	require.NoError(t, err)

	const workers = 1000
	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.Operation(w.ID, enums.WITHDRAW, 1)
			errCh <- err
		}()
	}
	wg.Wait()
	close(errCh)

	for e := range errCh {
		require.NoError(t, e)
	}
	got, err := repo.Get(w.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, got.Balance)
}
