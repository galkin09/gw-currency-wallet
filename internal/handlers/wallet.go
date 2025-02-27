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
// @Description  deposit user wallet
// @Tags         accounts
// @Param 		 Authorization header string true "JWT token"
// @Param		 amount body models.DepositReq true "Deposit query in json format"
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Router       /api/v1/wallet/deposit [post]
func (h *Handler) Deposit(c *gin.Context) {
	var dq storages.Deposit

	// Привязка JSON из запроса
	if err := c.ShouldBindJSON(&dq); err != nil {
		h.logger.Error("Could not bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Проверка обязательных полей
	if dq.Amount <= 0 || dq.Currency == "" {
		h.logger.Error("Invalid deposit request", zap.Any("request", dq))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "amount and currency are required"})
		return
	}

	// Получаем текущий баланс пользователя из базы данных
	username := c.GetString("username") // Предположим, что username сохранен в контексте
	wallet, err := h.storage.GetWalletByUsername(c, username)
	if err != nil {
		h.logger.Error("Failed to get wallet", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve wallet"})
		return
	}

	// Обновляем баланс
	switch dq.Currency {
	case "USD":
		wallet.Balance.USD += dq.Amount
	case "EUR":
		wallet.Balance.EUR += dq.Amount
	case "RUB":
		wallet.Balance.RUB += dq.Amount
	default:
		h.logger.Error("Unsupported currency", zap.String("currency", dq.Currency))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unsupported currency"})
		return
	}

	// Логируем новый баланс
	h.logger.Info("Updated wallet balance", zap.Any("balance", wallet.Balance))

	// Сохраняем обновленный баланс в базе данных
	if err := h.storage.Deposit(c, wallet); err != nil {
		h.logger.Error("Failed to deposit funds", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update wallet"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"message":     "Deposit successful",
		"new_balance": wallet.Balance,
	})
}

// Withdraw wallet with provided amount
// @Summary      Withdraw amount
// @Description  withdraw provided amount from user wallet
// @Tags         accounts
// @Param 		 Authorization header string true "JWT token"
// @Param		 amount body models.WithdrawReq true "Withdraw query in json format"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Проверка обязательных полей
	if wq.Amount <= 0 || wq.Currency == "" {
		h.logger.Error("Invalid withdraw request", zap.Any("request", wq))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "amount and currency are required"})
		return
	}

	// Получаем текущий баланс пользователя из базы данных
	username := c.GetString("username") // Предположим, что username сохранен в контексте
	wallet, err := h.storage.GetWalletByUsername(c, username)
	if err != nil {
		h.logger.Error("Failed to get wallet", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve wallet"})
		return
	}

	// Проверяем, достаточно ли средств для списания
	switch wq.Currency {
	case "USD":
		if wallet.Balance.USD < wq.Amount {
			h.logger.Error("Insufficient funds", zap.String("currency", wq.Currency))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
			return
		}
		wallet.Balance.USD -= wq.Amount
	case "EUR":
		if wallet.Balance.EUR < wq.Amount {
			h.logger.Error("Insufficient funds", zap.String("currency", wq.Currency))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
			return
		}
		wallet.Balance.EUR -= wq.Amount
	case "RUB":
		if wallet.Balance.RUB < wq.Amount {
			h.logger.Error("Insufficient funds", zap.String("currency", wq.Currency))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
			return
		}
		wallet.Balance.RUB -= wq.Amount
	default:
		h.logger.Error("Unsupported currency", zap.String("currency", wq.Currency))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unsupported currency"})
		return
	}

	// Логируем новый баланс
	h.logger.Info("Updated wallet balance", zap.Any("balance", wallet.Balance))

	// Сохраняем обновленный баланс в базе данных
	if err := h.storage.Withdraw(c, wallet); err != nil {
		h.logger.Error("Failed to withdraw funds", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update wallet"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"message":     "Withdrawal successful",
		"new_balance": wallet.Balance,
	})
}

// Exchange one currency to another with provided amount
//
//	@Summary      Exchanger endpoint
//	@Description  exchange one currency to another
//	@Tags         exchange
//	@Param 		 Authorization header string true "JWT token"
//	@Param		  amount body models.ExchangeReq true "Exchange query in json format"
//	@Accept       json
//	@Produce      json
//	@Success      200
//	@Failure      400
//	@Router       /api/v1/exchange [post]
//func (h *Handler) Exchange(c *gin.Context) {
//	var exchangeReq storages.Exchanger
//
//	username := c.GetString("username")
//	if username == "" {
//		h.logger.Error("Username not found in context")
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "username not found"})
//		return
//	}
//
//	if err := c.ShouldBindJSON(&exchangeReq); err != nil {
//		h.logger.Error("Could not bind JSON", zap.Error(err))
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
//		return
//	}
//
//	if exchangeReq.FromCurrency == "" || exchangeReq.ToCurrency == "" || exchangeReq.Amount <= 0 {
//		h.logger.Error("Invalid exchange request", zap.Any("request", exchangeReq))
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "from_currency, to_currency, and amount are required"})
//		return
//	}
//
//	// Получаем курс валют из кэша
//	rate, found := h.Cache.Get("rate")
//	if found {
//		h.logger.Debug("Rate fetched from cache", zap.Any("rate", rate))
//		exchangeReq.Rate = rate.(storages.Currency)
//	} else {
//		h.logger.Debug("Could not get rate from cache")
//		// Если курс не найден в кэше, можно запросить его из внешнего сервиса
//		// Например: exchangeReq.Rate, err = h.currencyService.GetRate(exchangeReq.FromCurrency, exchangeReq.ToCurrency)
//		// if err != nil { ... }
//	}
//
//	// Получаем кошелек пользователя
//	wallet, err := h.storage.GetWalletByUsername(c, username)
//	if err != nil {
//		h.logger.Error("Failed to get wallet", zap.Error(err))
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve wallet"})
//		return
//	}
//
//	// Выполняем обмен валют
//	updatedWallet, err := h.storage.Exchange(c, wallet, exchangeReq)
//	if err != nil {
//		h.logger.Error("Failed to exchange currency", zap.Error(err))
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange currency"})
//		return
//	}
//
//	// Логируем новый баланс
//	h.logger.Info("Updated wallet balance after exchange", zap.Any("balance", updatedWallet.Balance))
//
//	// Возвращаем успешный ответ
//	c.JSON(http.StatusOK, gin.H{
//		"message":     "Exchange successful",
//		"new_balance": updatedWallet.Balance,
//	})
//}

// ExchangeRates one currency to another with provided amount
//
//	@Summary      Exchanger endpoint
//	@Description  exchange one currency to another
//	@Tags         exchange
//	@Param 		 Authorization header string true "JWT token"
//	@Param		  amount body models.ExchangeReq true "Exchange query in json format"
//	@Accept       json
//	@Produce      json
//	@Success      200
//	@Failure      400
//	@Router       /api/v1/exchange [post]
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
