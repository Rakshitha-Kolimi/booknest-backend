package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"booknest/internal/domain"
	"booknest/internal/http/routes"
	"booknest/internal/middleware"
)

type userController struct {
	service domain.UserService
}

// NewUserController creates a new user controller instance
func NewUserController(service domain.UserService) domain.UserController {
	return &userController{service: service}
}

// RegisterRoutes registers all user routes
func (c *userController) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("")
	{
		auth.POST(routes.RegisterRoute, c.Register)
		auth.POST(routes.LoginRoute, c.Login)
		auth.POST(routes.ForgotPassword, c.ForgotPassword)
	}

	protected := r.Group("")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET(routes.UserRoute, c.GetUser)
		protected.DELETE(routes.UserRoute, c.DeleteUser)
		protected.POST("/verify-email", c.VerifyEmail)
		protected.POST("/verify-mobile", c.VerifyMobile)
		protected.POST("/resend-email-verification", c.ResendEmailVerification)
		protected.POST("/resend-mobile-otp", c.ResendMobileOTP)
		protected.POST("/reset-password", c.ResetPassword)
	}
}

// Register handles user registration
func (c *userController) Register(ctx *gin.Context) {
	var input domain.UserInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Register(ctx, input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please verify your email and mobile.",
	})
}

// Login handles user login
func (c *userController) Login(ctx *gin.Context) {
	var input domain.LoginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that either email or mobile is provided
	if input.Email == "" && input.Mobile == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email or mobile is required"})
		return
	}

	token, err := c.service.Login(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Login successful",
	})
}

// GetUser retrieves a user by ID
func (c *userController) GetUser(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := c.service.FindUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user account
func (c *userController) DeleteUser(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Verify that the user can only delete their own account
	userIDFromCtx, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if userIDFromCtx.(uuid.UUID) != id {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "you can only delete your own account"})
		return
	}

	if err := c.service.DeleteUser(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// ForgotPassword initiates password reset process
func (c *userController) ForgotPassword(ctx *gin.Context) {
	var input domain.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Email == "" && input.Mobile == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email or mobile is required"})
		return
	}

	// TODO: Implement forgot password logic with sending reset link/OTP
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password reset link has been sent to your email/mobile",
	})
}

// VerifyEmail verifies user email with token
func (c *userController) VerifyEmail(ctx *gin.Context) {
	var input struct {
		Token string `json:"token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.VerifyEmail(ctx, input.Token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
	})
}

// VerifyMobile verifies user mobile with OTP
func (c *userController) VerifyMobile(ctx *gin.Context) {
	var input struct {
		OTP string `json:"otp" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.VerifyMobile(ctx, input.OTP); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired OTP"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Mobile verified successfully",
	})
}

// ResendEmailVerification resends email verification token
func (c *userController) ResendEmailVerification(ctx *gin.Context) {
	userIDFromCtx, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDFromCtx.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id format"})
		return
	}

	if err := c.service.ResendEmailVerification(ctx, userID); err != nil {
		if errors.Is(err, errors.New("email already verified")) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "email already verified"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent successfully",
	})
}

// ResendMobileOTP resends mobile OTP
func (c *userController) ResendMobileOTP(ctx *gin.Context) {
	userIDFromCtx, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDFromCtx.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id format"})
		return
	}

	if err := c.service.ResendMobileOTP(ctx, userID); err != nil {
		if errors.Is(err, errors.New("mobile already verified")) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "mobile already verified"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Mobile OTP sent successfully",
	})
}

// ResetPassword resets user password
func (c *userController) ResetPassword(ctx *gin.Context) {
	userIDFromCtx, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDFromCtx.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id format"})
		return
	}

	var input struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.ResetPassword(ctx, userID, input.NewPassword); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}
