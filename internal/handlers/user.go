package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gw-currency-wallet/internal/auth"
	"gw-currency-wallet/internal/storages"
	"net/http"
	"time"
)

// RegisterUser adds new user account
//
//	@Summary      Register user
//	@Description  Регистрация нового пользователя
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Param input body storages.User true "user info"
//	@Success      201
//	@Failure      400
//	@Router       /api/v1/register [post]
func (h *Handler) RegisterUser(ctx *gin.Context) {
	var user storages.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		h.logger.Error("Error binding JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user == (storages.User{}) {
		h.logger.Error("User is nil", zap.Any("user", user))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User is nil"})
		return
	}

	if err := h.storage.RegisterUser(ctx, user); err != nil {
		h.logger.Error("Username or email already exists", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Successfully registered user", zap.Any("user", user))
	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// LoginUser authorizes  adds new user account
//
//		@Summary      Authorize  user
//		@Description  Авторизация пользователя
//		@Tags         users
//		@Accept       json
//		@Produce      json
//	    @Param input body storages.User true "Данные пользователя"
//		@Success      200
//		@Failure      400
//		@Router       /api/v1/login [post]
func (h *Handler) LoginUser(ctx *gin.Context) {
	var user storages.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		h.logger.Error("Error binding JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtToken, err := auth.GenerateToken(user, 10*time.Minute)
	if err != nil {
		h.logger.Error("login Error", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Successfully logged in", zap.Any("jwtToken", jwtToken))
	ctx.JSON(http.StatusOK, gin.H{"token": jwtToken})
}

// GetBalance godoc
//
//	@Summary      Shows wallet balance
//	@Description  Показывает баланс пользователя на счёте по юзернейму
//	@Tags         users, wallets
//	@Param 		  Authorization header string true "JWT token"
//	@Accept       json
//	@Produce      json
//	@Success      200
//	@Failure      400
//	@Router       /api/v1/balance [get]
func (h *Handler) GetBalance(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		h.logger.Error("Username not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	usernameStr, ok := username.(string)
	if !ok {
		h.logger.Error("Invalid username format in context")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid username format"})
		return
	}

	user := storages.User{
		Username: usernameStr,
	}

	h.logger.Info("Getting balance for user", zap.String("username", user.Username))

	balance, err := h.storage.GetBalance(ctx, user)
	if err != nil {
		h.logger.Error("GetBalance Error", zap.String("username", user.Username), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Successfully got balance", zap.Any("balance", balance))
	ctx.JSON(http.StatusOK, gin.H{"balance": balance})
}
