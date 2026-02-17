package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"booknest/internal/domain"
)

func getUserID(ctx *gin.Context) (uuid.UUID, error) {
	raw, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("missing user_id in context")
	}

	switch v := raw.(type) {
	case string:
		return uuid.Parse(v)
	case uuid.UUID:
		return v, nil
	default:
		return uuid.Nil, errors.New("invalid user_id type")
	}
}

func getUserRole(ctx *gin.Context) (domain.UserRole, error) {
	raw, exists := ctx.Get("user_role")
	if !exists {
		return "", errors.New("missing user_role in context")
	}

	switch v := raw.(type) {
	case string:
		return domain.UserRole(v), nil
	case domain.UserRole:
		return v, nil
	default:
		return "", errors.New("invalid user_role type")
	}
}
