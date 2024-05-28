package handler

import (
	"expenses_tracker/internal/model"
	"expenses_tracker/internal/pkg/auth"
	"expenses_tracker/internal/pkg/jwt"
	"expenses_tracker/internal/pkg/password"
	"expenses_tracker/internal/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	jwtService     *jwt.JwtService
	userRepository repository.UserRepository
}

func RegisterUserRoutes(router *gin.Engine, jwtService *jwt.JwtService, userRepository repository.UserRepository) {
	handler := userHandler{
		jwtService:     jwtService,
		userRepository: userRepository,
	}

	userRouterGroup := router.Group("/user")

	userRouterGroup.POST("/register", handler.register)
	userRouterGroup.GET("/login", handler.login)
	userRouterGroup.GET("/", handler.get).Use(auth.GetAuthMiddleware(jwtService))
}

func (h *userHandler) register(c *gin.Context) {
	var user model.UserModel
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user object"})
		return
	}

	hashedPassword, err := password.HashPassword(user.PasswordHash)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.PasswordHash = hashedPassword
	err = h.userRepository.Create(user)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("user", user)
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *userHandler) login(c *gin.Context) {
	login, passwordParam := c.Query("login"), c.Query("password")
	if login == "" || passwordParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide login and password"})
		return
	}

	user, err := h.userRepository.FindByLogin(login)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "hi": err.Error()})
		return
	}

	if !password.ComparePassword(user.PasswordHash, passwordParam) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	token, err := h.jwtService.GenerateToken(user.Id)
	if err != nil {
		print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *userHandler) get(c *gin.Context) {
	userId, ok := auth.GetUserId(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userRepository.FindById(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
