package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"booknest/internal/domain"
)

type UserController struct {
	service domain.UserService
}

func NewUserController(s domain.UserService) *UserController {
	return &UserController{service: s}
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user and returns a JWT token
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			input	body		domain.UserInput	true	"User input"
//	@Success		201		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Router			/register [post]
func (c *UserController) RegisterUser(ctx *gin.Context) {
	var input domain.UserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.RegisterUser(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "registration successful",
		"token":   token,
	})
}

// LoginUser godoc
//
//	@Summary		Login user
//	@Description	Authenticate a user and return a JWT token
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			input	body	map[string]string	true	"Login credentials"
//	@Success		200		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Router			/login [post]
func (c *UserController) LoginUser(ctx *gin.Context) {
	var creds map[string]string
	if err := ctx.ShouldBindJSON(&creds); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, password := creds["email"], creds["password"]
	if email == "" || password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	token, err := c.service.LoginUser(ctx, email, password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}

// GetUsers godoc
//
//	@Summary		Get all users
//	@Tags			User
//	@Produce		json
//	@Success		200	{array}	domain.User
//	@Router			/users [get]
func (c *UserController) GetUsers(ctx *gin.Context) {
	users, err := c.service.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// GetUserByID godoc
//
//	@Summary		Get user by ID
//	@Tags			User
//	@Produce		json
//	@Param			id	path	string	true	"User ID"
//	@Success		200	{object}	domain.User
//	@Failure		404	{object}	map[string]string
//	@Router			/users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := c.service.GetUserByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// DeleteUser godoc
//
//	@Summary		Delete user by ID (soft delete)
//	@Tags			User
//	@Produce		json
//	@Param			id	path	string	true	"User ID"
//	@Success		200	{object}	map[string]string
//	@Router			/users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := c.service.DeleteUser(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
