package handler

import (
	"expenses_tracker/internal/model"
	"expenses_tracker/internal/pkg/auth"
	"expenses_tracker/internal/pkg/jwt"
	"expenses_tracker/internal/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	transactionRepository repository.TransactionRepository
}

func RegisterTransactionRoutes(router *gin.Engine, jwtService *jwt.JwtService, transactionRepositor repository.TransactionRepository) {
	handler := transactionHandler{
		transactionRepository: transactionRepositor,
	}

	transactionRouterGroup := router.Group("/transaction").Use(auth.GetAuthMiddleware(jwtService))

	transactionRouterGroup.POST("", handler.create)
	transactionRouterGroup.GET("", handler.get)
	transactionRouterGroup.DELETE("", handler.deleteTransaction)
}

func (h *transactionHandler) create(c *gin.Context) {
	userId, ok := auth.GetUserId(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var transaction model.Transaction
	if err := c.BindJSON(&transaction); err != nil || transaction.Price == 0 || transaction.CategoryId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction object"})
		return
	}

	transaction.UserId = userId
	err := h.transactionRepository.CreateTransaction(transaction)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "can't create transaction",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *transactionHandler) get(c *gin.Context) {
	userId, ok := auth.GetUserId(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type GetTransactionsInput struct {
		CategoryIds []int64               `json:"categoryIds" binding:"required"`
		Pagination  repository.Pagination `json:"pagination" binding:"required"`
	}

	categoryIdsParam := c.Query("categoryIds")
	var categoryIds []int64
	if categoryIdsParam != "" {
		categoryIdsStr := strings.Split(categoryIdsParam, ",")
		for _, idStr := range categoryIdsStr {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid categoryId: " + idStr})
				return
			}
			categoryIds = append(categoryIds, id)
		}
	}

	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing page parameter"})
		return
	}

	items, err := strconv.ParseInt(c.Query("items"), 10, 64)
	if err != nil || items <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing items parameter"})
		return
	}

	pagination := repository.Pagination{
		Page:  page,
		Items: items,
	}

	resolvedPagination := repository.ResolvePagination(&pagination)

	transactions, err := h.transactionRepository.GetTransactions(userId, categoryIds, resolvedPagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *transactionHandler) deleteTransaction(c *gin.Context) {
	userId, ok := auth.GetUserId(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type DeleteTransactionInput struct {
		TransactionId int64 `json:"id" binding:"required"`
	}

	var input DeleteTransactionInput
	if err := c.BindJSON(&input); err != nil || input.TransactionId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input object"})
		return
	}

	transaction, err := h.transactionRepository.GetTransactionById(input.TransactionId)
	if err != nil || transaction.UserId != userId {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	err = h.transactionRepository.DeleteTransaction(transaction.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}

	c.String(http.StatusOK, "OK")
}
