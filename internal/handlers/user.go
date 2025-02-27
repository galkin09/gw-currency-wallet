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
//	@Description  register new users
//	@Tags         accounts
//	@Accept       json
//	@Produce      json
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
//	@Summary      Authorize existing user
//	@Description  authorize users
//	@Tags         accounts
//	@Accept       json
//	@Produce      json
//	@Success      200
//	@Failure      400
//	@Router       /api/v1/login [post]
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

// GetBalance returns wallet balance
//
//	@Summary      Shows wallet balance
//	@Description  shows user wallet balance
//	@Tags         accounts, wallets
//	@Param 		  Authorization header string true "JWT token"
//	@Accept       json
//	@Produce      json
//	@Success      200
//	@Failure      400
//	@Router       /api/v1/balance [get]
func (h *Handler) GetBalance(ctx *gin.Context) {
	var user storages.User

	user.Username = ctx.Param("username")
	h.logger.Info("Getting balance for user", zap.Any("user", user))

	balance, err := h.storage.GetBalance(ctx, user)
	if err != nil {
		h.logger.Error("GetBalance Error", zap.Any("user", user), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Successfully got balance for wallet", zap.Any("wallet", storages.Wallet{}))
	ctx.JSON(http.StatusOK, gin.H{"balance": balance})
}
