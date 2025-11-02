package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saifwork/mock-service/internal/api/responses"
	"github.com/saifwork/mock-service/internal/core/config"
	"github.com/saifwork/mock-service/internal/dtos"
	"github.com/saifwork/mock-service/internal/middlewares"
	"github.com/saifwork/mock-service/internal/services"
)

// AuthHandler handles user authentication routes.
type AuthHandler struct {
	service *services.AuthService
	cfg     *config.Config
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(service *services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{service: service, cfg: cfg}
}

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/signup", h.Signup)
		authRoutes.POST("/login", h.Login)
		authRoutes.GET("/verify-email", h.VerifyEmail)
		authRoutes.POST("/forgot-password", h.ForgotPassword)
		authRoutes.POST("/reset-forgot-password", h.ResetForgotPassword)

		// ✅ Protected route
		authRoutes.POST("/reset-password", middlewares.AuthMiddleware(h.cfg), h.ResetPassword)
		authRoutes.GET("/me", middlewares.AuthMiddleware(h.cfg), h.GetCurrentUser)
	}
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		responses.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.service.GetUserByID(userID.(string))
	if err != nil {
		responses.JSONError(c, http.StatusNotFound, "User not found")
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "User fetched successfully", user)
}

// -------------------- Signup --------------------
func (h *AuthHandler) Signup(c *gin.Context) {
	var req dtos.UserRegisterRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.JSONError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.Signup(&req); err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusCreated, "User registered successfully, please verify your email", nil)
}

// -------------------- Verify Email --------------------
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		responses.JSONError(c, http.StatusBadRequest, "Missing verification token")
		return
	}

	if err := h.service.VerifyEmail(token); err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Email verified successfully", nil)
}

// -------------------- Login --------------------
func (h *AuthHandler) Login(c *gin.Context) {
	var req dtos.UserLoginRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.JSONError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, err := h.service.Login(&req)
	if err != nil {
		responses.JSONError(c, http.StatusUnauthorized, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Login successful", gin.H{"token": token})
}

// -------------------- Forgot Password --------------------
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var payload struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		responses.JSONError(c, http.StatusBadRequest, "Invalid email format")
		return
	}

	if err := h.service.ForgotPassword(payload.Email); err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Password reset link sent to your email", nil)
}

// -------------------- Reset Forgot Password --------------------
func (h *AuthHandler) ResetForgotPassword(c *gin.Context) {
	var payload struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		responses.JSONError(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	if err := h.service.ResetForgotPassword(payload.Token, payload.NewPassword); err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Password has been reset successfully", nil)
}

// -------------------- Reset Password (Logged-in User) --------------------
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var payload struct {
		OldPwd string `json:"oldPassword" binding:"required"`
		NewPwd string `json:"newPassword" binding:"required,min=6"`
	}

	// 🔹 Parse JSON body
	if err := c.ShouldBindJSON(&payload); err != nil {
		responses.JSONError(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	// 🔹 Get userID from JWT (middleware must set this)
	userID, exists := c.Get("userID")
	if !exists {
		responses.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 🔹 Call service with the authenticated userID
	if err := h.service.ResetPassword(userID.(string), payload.OldPwd, payload.NewPwd); err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Password changed successfully", nil)
}
