package main

import (
	"itk-academy-test/config"
	"itk-academy-test/internal/handlers"
	"itk-academy-test/internal/models"
	"itk-academy-test/internal/repository"
	"itk-academy-test/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	postgresConfig := config.PostgresConfig{}
	postgresConfig = postgresConfig.Load()

	db, err := gorm.Open(postgres.Open(postgresConfig.Print()))
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = db.AutoMigrate(
		&models.Wallet{},
	)

	if err != nil {
		log.Fatal("Failed to migrate the database", err)
	}

	walletRepository := &repository.WalletGORMRepository{DB: db}
	walletService := services.New(walletRepository)
	walletHandler := handlers.New(walletService)

	walletHandler.Initialize(r)

	r.Run(":9090")
}
