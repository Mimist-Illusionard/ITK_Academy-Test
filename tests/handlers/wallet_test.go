package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	enums "itk-academy-test/internal"
	"itk-academy-test/internal/dto"
	"itk-academy-test/internal/handlers"
	"itk-academy-test/internal/models"
	"itk-academy-test/internal/repository"
	"itk-academy-test/internal/services"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
	return db
}

func newRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := newDB(t)
	repo := &repository.WalletGORMRepository{DB: db}
	svc := services.New(repo)
	h := handlers.New(svc)

	r := gin.New()

	h.Initialize(r)
	return r
}

func TestCreateWallet(t *testing.T) {
	r := newRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/wallets/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.WalletResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotZero(t, resp.WalletID)
	assert.Equal(t, 0, resp.Balance)
	assert.Equal(t, "Wallet created", resp.Message)
}

func TestAmount(t *testing.T) {
	r := newRouter(t)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/wallets/", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	var created dto.WalletResponse
	err := json.Unmarshal(w1.Body.Bytes(), &created)
	assert.NoError(t, err)
	assert.NotZero(t, created.WalletID)
	assert.Equal(t, 0, created.Balance)

	opBody, _ := json.Marshal(dto.WalletOperationRequest{
		WalletID:      created.WalletID,
		OperationType: string(enums.DEPOSIT),
		Amount:        150,
	})
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/wallet/", bytes.NewReader(opBody))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/wallets/"+fmt.Sprint(created.WalletID), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.WalletResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, created.WalletID, resp.WalletID)
	assert.Equal(t, 150, resp.Balance)
	assert.Equal(t, "", resp.Message)
}

func TestOperation_Deposit(t *testing.T) {
	r := newRouter(t)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/wallets/", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	var created dto.WalletResponse
	_ = json.Unmarshal(w1.Body.Bytes(), &created)

	body, _ := json.Marshal(dto.WalletOperationRequest{
		WalletID:      created.WalletID,
		OperationType: string(enums.DEPOSIT),
		Amount:        200,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/wallet/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.WalletResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, created.WalletID, resp.WalletID)
	assert.Equal(t, 200, resp.Balance)
	assert.Equal(t, "Operation completed successfully", resp.Message)
}

func TestOperation_Withdraw_Insufficient(t *testing.T) {
	r := newRouter(t)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/wallets/", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	var created dto.WalletResponse
	_ = json.Unmarshal(w1.Body.Bytes(), &created)

	body, _ := json.Marshal(dto.WalletOperationRequest{
		WalletID:      created.WalletID,
		OperationType: string(enums.WITHDRAW),
		Amount:        50,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/wallet/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Insufficient funds")
}

func TestDeleteWallet(t *testing.T) {
	r := newRouter(t)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/wallets/", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	var created dto.WalletResponse
	_ = json.Unmarshal(w1.Body.Bytes(), &created)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/wallets/"+fmt.Sprint(created.WalletID), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Wallet deleted")
}
