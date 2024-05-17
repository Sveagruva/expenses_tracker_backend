package main

import (
	"expenses_tracker/internal/config"
	"expenses_tracker/internal/handler"
	"expenses_tracker/internal/pkg/jwt"
	"expenses_tracker/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfigFromEnv(".env")
	db, err := repository.GetSqliteDb(cfg.DB.Path)
	if err != nil {
		panic(err)
	}

	jwtService := &jwt.JwtService{PrivateKey: cfg.Jwt.PrivateKey}
	userRepo := repository.GetUserRepository(db)
	transactionRepo := repository.GetTransactionRepository(db)
	transactionCategoryRepo := repository.GetTransactionCategoryRepository(db)

	router := gin.Default()

	handler.RegisterUserRoutes(router, jwtService, userRepo)
	handler.RegisterTransactionRoutes(router, jwtService, transactionRepo)
	handler.RegisterTransactionCategoryRoutes(router, jwtService, transactionCategoryRepo)

	router.Run()
}
