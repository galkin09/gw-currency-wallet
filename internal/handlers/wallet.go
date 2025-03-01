package handlers

import (
	pb "github.com/galkin09/proto-exchange/exchange"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gw-currency-wallet/internal/storages"
	"net/http"
)

// Deposit wallet with provided amount
//
// @Summary      Deposit balance
// @Description  Пополнить баланс пользователя
// @Tags         wallets, users
// @Param 		 Authorization header string true "JWT token"
// @Param		 amount body storages.Deposit true "Deposit query in json format"
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Router       /api/v1/wallet/deposit [post]
func (h *Handler) Deposit(c *gin.Context) {
	var dq storages.Deposit
	var wallet storages.Wallet

	if err := c.ShouldBindJSON(&dq); err != nil {
		h.logger.Error("Could not bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		h.logger.Error("Username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	wallet, err := h.storage.GetWalletByUsername(c, username.(string))
	if err != nil {
		h.logger.Error("Could not get wallet", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve wallet"})
		return
	}

	switch dq.Currency {
	case "RUB":
		wallet.Balance.RUB += dq.Amount
	case "USD":
		wallet.Balance.USD += dq.Amount
	case "EUR":
		wallet.Balance.EUR += dq.Amount
	default:
		h.logger.Error("Unsupported currency", zap.String("currency", dq.Currency))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unsupported currency"})
		return
	}

	h.logger.Info("Updated wallet balance", zap.Any("balance", wallet.Balance))

	if err := h.storage.Deposit(c, wallet, dq.Currency, dq.Amount); err != nil {
		h.logger.Error("Could not deposit", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Account topped up successfully",
		"new_balance": wallet.Balance,
	})
}

// Withdraw wallet with provided amount
// @Summary      Withdraw amount
// @Description  Снять средства со счёта пользователя
// @Tags         users, wallets
// @Param 		 Authorization header string true "JWT token"
// @Param		 amount body storages.Withdraw true "Withdraw query in json format"
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Router       /api/v1/wallet/withdraw [post]
func (h *Handler) Withdraw(c *gin.Context) {
	var wq storages.Withdraw

	// Привязка JSON из запроса
	if err := c.ShouldBindJSON(&wq); err != nil {
		h.logger.Error("Could not bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Debug("Withdraw request", zap.Any("request", wq))

	username, exists := c.Get("username")
	if !exists {
		h.logger.Error("Username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	wallet, err := h.storage.GetWalletByUsername(c, username.(string)) // Предположим, что wq.WalletID содержит ID кошелька
	if err != nil {
		h.logger.Error("Could not get wallet", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve wallet"})
		return
	}

	switch wq.Currency {
	case "USD":
		if wallet.Balance.USD < wq.Amount {
			h.logger.Error("Insufficient funds", zap.Float64("amount", wq.Amount), zap.Float64("balance", wallet.Balance.USD))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
			return
		}
		wallet.Balance.USD -= wq.Amount
	case "EUR":
		if wallet.Balance.EUR < wq.Amount {
			h.logger.Error("Insufficient funds", zap.Float64("amount", wq.Amount), zap.Float64("balance", wallet.Balance.EUR))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
			return
		}
		wallet.Balance.EUR -= wq.Amount
	case "RUB":
		if wallet.Balance.RUB < wq.Amount {
			h.logger.Error("Insufficient funds", zap.Float64("amount", wq.Amount), zap.Float64("balance", wallet.Balance.RUB))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
			return
		}
		wallet.Balance.RUB -= wq.Amount
	default:
		h.logger.Error("Unsupported currency", zap.String("currency", wq.Currency))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unsupported currency"})
		return
	}

	h.logger.Debug("Updated wallet balance", zap.Any("balance", wallet.Balance))

	if err := h.storage.Withdraw(c, wallet, wq.Currency, wq.Amount); err != nil {
		h.logger.Error("Could not withdraw", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Funds withdrawn successfully",
		"new_balance": wallet.Balance,
	})
}

// Exchange one currency to another with provided amount
//
//	@Summary      Exchanger endpoint
//	@Description  Позволяет обменять валюту на другую, курс можно узнать в /api/v1/exchange/rates
//	@Tags         exchange
//	@Param 		  Authorization header string true "JWT token"
//	@Param		  amount body storages.Exchanger true "Exchange query in json format"
//	@Accept       json
//	@Produce      json
//	@Success      200
//	@Failure      400
//	@Router       /api/v1/exchange [post]
func (h *Handler) Exchange(c *gin.Context) {
	var ex storages.Exchanger

	if err := c.ShouldBindJSON(&ex); err != nil {
		h.logger.Error("Could not bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Debug("Exchange request", zap.Any("request", ex))

	username, exists := c.Get("username")
	if !exists {
		h.logger.Error("Username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	wallet, err := h.storage.GetWalletByUsername(c, username.(string)) // Предположим, что ex.WalletID содержит ID кошелька
	if err != nil {
		h.logger.Error("Could not get wallet", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve wallet"})
		return
	}

	rate, err := h.exch.GetExchangeRateForCurrency(c.Request.Context(), &pb.CurrencyRequest{
		FromCurrency: ex.FromCurrency,
		ToCurrency:   ex.ToCurrency,
	})
	if err != nil {
		h.logger.Error("Could not get exchange rate", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get exchange rate"})
		return
	}

	convertedAmount := ex.Amount * float64(rate.Rate)

	switch ex.FromCurrency {
	case "USD":
		if wallet.Balance.USD < ex.Amount {
			h.logger.Error("Insufficient funds", zap.Float64("amount", ex.Amount), zap.Float64("balance", wallet.Balance.USD))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
			return
		}
		wallet.Balance.USD -= ex.Amount
	case "EUR":
		if wallet.Balance.EUR < ex.Amount {
			h.logger.Error("Insufficient funds", zap.Float64("amount", ex.Amount), zap.Float64("balance", wallet.Balance.EUR))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
			return
		}
		wallet.Balance.EUR -= ex.Amount
	case "RUB":
		if wallet.Balance.RUB < ex.Amount {
			h.logger.Error("Insufficient funds", zap.Float64("amount", ex.Amount), zap.Float64("balance", wallet.Balance.RUB))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
			return
		}
		wallet.Balance.RUB -= ex.Amount
	default:
		h.logger.Error("Unsupported currency", zap.String("currency", ex.FromCurrency))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unsupported currency"})
		return
	}

	switch ex.ToCurrency {
	case "USD":
		wallet.Balance.USD += convertedAmount
	case "EUR":
		wallet.Balance.EUR += convertedAmount
	case "RUB":
		wallet.Balance.RUB += convertedAmount
	default:
		h.logger.Error("Unsupported currency", zap.String("currency", ex.ToCurrency))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unsupported currency"})
		return
	}

	h.logger.Debug("Updated wallet balance", zap.Any("balance", wallet.Balance))

	if err := h.storage.UpdateWalletBalance(c, wallet); err != nil {
		h.logger.Error("Could not update wallet balance", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Exchange completed successfully",
		"exchanged_amount": convertedAmount,
		"new_balance":      wallet.Balance,
	})
}

// ExchangeRates all rates in exchanger
//
//	@Summary      Exchanger endpoint
//	@Description  Позволяет узнать актуальный курс по отношению к доллару
//	@Tags         exchange
//	@Param 		  Authorization header string true "JWT token"
//	@Accept       json
//	@Produce      json
//	@Success      200
//	@Failure      400
//	@Router       /api/v1/exchange/rates [get]
func (h *Handler) ExchangeRates(c *gin.Context) {
	var resp storages.Rates

	eResp, err := h.exch.GetExchangeRates(c, &pb.Empty{})
	if err != nil {
		// todo: get rates from cache
		h.logger.Error("Could not get exchange rates", zap.Error(err))
	}

	if eResp == nil {
		h.logger.Error("No exchange rates found", zap.Error(err))
		return
	}

	mapRates := eResp.GetRates()
	resp.USD = float64(mapRates["USD"])
	resp.EUR = float64(mapRates["EUR"])
	resp.RUB = float64(mapRates["RUB"])

	c.JSON(http.StatusOK, resp)
}
