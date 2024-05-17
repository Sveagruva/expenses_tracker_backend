package handler

import (
	"expenses_tracker/internal/model"
	"expenses_tracker/internal/pkg/auth"
	"expenses_tracker/internal/pkg/jwt"
	"expenses_tracker/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionCategoryHandler struct {
	transactionCategoryRepository repository.TransactionCategoryRepository
}

func RegisterTransactionCategoryRoutes(router *gin.Engine, jwtService *jwt.JwtService, transactionCategoryRepository repository.TransactionCategoryRepository) {
	handler := transactionCategoryHandler{
		transactionCategoryRepository: transactionCategoryRepository,
	}

	transactionCategoryRouterGroup := router.Group("/transaction/category").Use(auth.GetAuthMiddleware(jwtService))

	transactionCategoryRouterGroup.POST("", handler.create)
	transactionCategoryRouterGroup.GET("", handler.get)
	transactionCategoryRouterGroup.DELETE("", handler.deleteCategory)
}

func (h *transactionCategoryHandler) create(c *gin.Context) {
	userId, ok := auth.GetUserId(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var category model.TransactionCategory
	if err := c.BindJSON(&category); err != nil || category.Name == "" || category.Color == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category object"})
		return
	}

	category.UserId = userId
	err := h.transactionCategoryRepository.CreateTransactionCategory(category)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "can't create category",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *transactionCategoryHandler) get(c *gin.Context) {
	userId, ok := auth.GetUserId(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	categories, err := h.transactionCategoryRepository.GetTransactionCategories(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *transactionCategoryHandler) deleteCategory(c *gin.Context) {
	userId, ok := auth.GetUserId(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type DeleteCategoryInput struct {
		CategoryId int64 `json:"id" binding:"required"`
	}

	var input DeleteCategoryInput
	if err := c.BindJSON(&input); err != nil || input.CategoryId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input object"})
		return
	}

	category, err := h.transactionCategoryRepository.GetTransactionCategoryById(input.CategoryId)
	if err != nil || category.UserId != userId {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	err = h.transactionCategoryRepository.DeleteTransactionCategory(category.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}

	c.String(http.StatusOK, "OK")
}
