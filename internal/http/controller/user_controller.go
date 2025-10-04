package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"booknest/internal/domain"
)

var dummyUsers = []domain.UserInput{
	{FirstName: "Alice", LastName: "Johnson", Email: "alice.johnson@example.com", Password: "pass1234"},
	{FirstName: "Brian", LastName: "Smith", Email: "brian.smith@example.com", Password: "brian987"},
	{FirstName: "Catherine", LastName: "Williams", Email: "catherine.williams@example.com", Password: "catpass1"},
	{FirstName: "David", LastName: "Brown", Email: "david.brown@example.com", Password: "david321"},
	{FirstName: "Ella", LastName: "Davis", Email: "ella.davis@example.com", Password: "ella4567"},
}

func GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dummyUsers)
}

func AddUser(ctx *gin.Context) {
	// Bind the input
	var in domain.UserInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// add the expense
	input := domain.UserInput{
		FirstName: in.FirstName,
		LastName: in.LastName,
		Email: in.Email,
		Password: in.Password,
	}
	dummyUsers = append(dummyUsers, input)

	// Return the data
	ctx.JSON(http.StatusOK, input)
}
