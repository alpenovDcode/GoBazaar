package handlers

import (
	"net/http"
	"strconv"

	"github.com/alpewa/GoBazaar/internal/auth/service"
	"github.com/alpewa/GoBazaar/internal/common/models"
	"github.com/gin-gonic/gin"
)

// AuthHandler обрабатывает HTTP запросы аутентификации
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler создает новый auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// ErrorResponse представляет структуру ответа с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// SuccessResponse представляет структуру успешного ответа
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse представляет структуру ответа с пагинацией
type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

// Register обрабатывает регистрацию пользователя
// @Summary Регистрация нового пользователя
// @Description Регистрирует новый аккаунт пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Запрос на регистрацию"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		switch err {
		case service.ErrValidationFailed:
			statusCode = http.StatusBadRequest
			errorCode = "validation_failed"
		case service.ErrUserExists:
			statusCode = http.StatusConflict
			errorCode = "user_exists"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login обрабатывает вход пользователя
// @Summary Вход пользователя
// @Description Аутентифицирует пользователя и возвращает токены
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Запрос на вход"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		switch err {
		case service.ErrInvalidCredentials:
			statusCode = http.StatusUnauthorized
			errorCode = "invalid_credentials"
		case service.ErrUserNotActive:
			statusCode = http.StatusUnauthorized
			errorCode = "user_not_active"
		case service.ErrValidationFailed:
			statusCode = http.StatusBadRequest
			errorCode = "validation_failed"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken обрабатывает обновление токена
// @Summary Обновление access токена
// @Description Обновляет access токен используя refresh токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Запрос на обновление токена"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		statusCode := http.StatusUnauthorized
		errorCode := "invalid_token"

		c.JSON(statusCode, ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout обрабатывает выход пользователя
// @Summary Выход пользователя
// @Description Отзывает refresh токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Запрос на выход"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Successfully logged out",
	})
}

// GetProfile получает профиль пользователя
// @Summary Получение профиля пользователя
// @Description Получает информацию о текущем пользователе
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer access_token"
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	response, err := h.authService.GetUser(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProfile обновляет профиль пользователя
// @Summary Обновление профиля пользователя
// @Description Обновляет информацию о пользователе
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer access_token"
// @Param request body map[string]interface{} true "Данные для обновления"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := h.authService.UpdateUser(c.Request.Context(), userID.(uint), updates)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		if err == service.ErrUserExists {
			statusCode = http.StatusConflict
			errorCode = "email_exists"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ChangePassword изменяет пароль пользователя
// @Summary Изменение пароля
// @Description Изменяет пароль пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer access_token"
// @Param request body models.ChangePasswordRequest true "Запрос на изменение пароля"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), userID.(uint), &req); err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		switch err {
		case service.ErrInvalidCredentials:
			statusCode = http.StatusUnauthorized
			errorCode = "invalid_current_password"
		case service.ErrValidationFailed:
			statusCode = http.StatusBadRequest
			errorCode = "validation_failed"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Password changed successfully",
	})
}

// GetUsers получает список всех пользователей (для администраторов)
// @Summary Получение списка пользователей
// @Description Получает список всех пользователей с пагинацией
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer access_token"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Success 200 {object} PaginatedResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /auth/users [get]
func (h *AuthHandler) GetUsers(c *gin.Context) {
	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	users, total, err := h.authService.GetAllUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:  users,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// SearchUsers ищет пользователей по имени
// @Summary Поиск пользователей
// @Description Ищет пользователей по имени или фамилии
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer access_token"
// @Param q query string true "Поисковый запрос"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/users/search [get]
func (h *AuthHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Search query is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	users, total, err := h.authService.SearchUsers(c.Request.Context(), query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:  users,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// Health проверка здоровья сервиса
// @Summary Проверка здоровья
// @Description Возвращает статус здоровья auth сервиса
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Router /auth/health [get]
func (h *AuthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Auth service is healthy",
		Data: map[string]interface{}{
			"status":    "ok",
			"timestamp": gin.H{"time": "now"},
		},
	})
}

// ValidateToken валидирует токен (внутренний метод)
func (h *AuthHandler) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	return h.authService.ValidateAccessToken(tokenString)
}
