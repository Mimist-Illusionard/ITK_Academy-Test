package tests

import (
	"context"
	"fmt"
	enums "itk-academy-test/internal"
	"itk-academy-test/internal/models"
	"itk-academy-test/internal/repository"
	"itk-academy-test/internal/services"
	"log"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	epg *embeddedpostgres.EmbeddedPostgres
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := startEmbeddedPostgres(ctx); err != nil {
		log.Fatalf("Failed to start embedded PostgreSQL: %v", err)
	}

	defer func() {
		if err := stopEmbeddedPostgres(); err != nil {
			log.Printf("Failed to stop embedded PostgreSQL: %v", err)
		}
	}()

	code := m.Run()
	os.Exit(code)
}

func startEmbeddedPostgres(ctx context.Context) error {
	if epg != nil {
		return nil
	}

	log.Println("Starting embedded PostgreSQL...")

	epg = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(5433).
		Username("postgres").
		Password("postgres").
		Database("test_db").
		RuntimePath("tmp_embedded_pg").
		Version(embeddedpostgres.V15))

	if err := epg.Start(); err != nil {
		return fmt.Errorf("failed to start embedded postgres: %w", err)
	}

	log.Println("PostgreSQL started, waiting for port to be ready...")

	if err := waitForPort(ctx, 5433, 30*time.Second); err != nil {
		return fmt.Errorf("postgres port not ready: %w", err)
	}

	if err := testConnection(ctx); err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	log.Println("Embedded PostgreSQL is ready!")
	return nil
}

func stopEmbeddedPostgres() error {
	if epg != nil {
		log.Println("Stopping embedded PostgreSQL...")
		if err := epg.Stop(); err != nil {
			return fmt.Errorf("failed to stop embedded postgres: %w", err)
		}
		epg = nil
		log.Println("Embedded PostgreSQL stopped")
	}
	return nil
}

func waitForPort(ctx context.Context, port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	address := fmt.Sprintf("localhost:%d", port)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for port %d", port)
			}

			conn, err := net.DialTimeout("tcp", address, 1*time.Second)
			if err == nil {
				conn.Close()
				return nil
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}

func testConnection(ctx context.Context) error {
	dsn := "host=localhost port=5433 user=postgres password=postgres dbname=test_db sslmode=disable client_encoding=UTF8"

	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				sqlDB, _ := db.DB()
				sqlDB.Close()
				return nil
			}
			time.Sleep(1 * time.Second)
		}
	}

	return fmt.Errorf("failed to connect after retries")
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "host=localhost port=5433 user=postgres password=postgres dbname=test_db sslmode=disable client_encoding=UTF8"

	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		t.Logf("Attempt %d: failed to connect, retrying...", i+1)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		t.Fatalf("failed to connect database after retries: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get SQL DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(2 * time.Hour)

	err = db.AutoMigrate(&models.Wallet{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func newDB(t *testing.T) *gorm.DB {
	t.Helper()
	return setupTestDB(t)
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

	_, err = svc.Operation(w.ID, enums.DEPOSIT, 500)
	require.NoError(t, err)

	const workers = 500
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
