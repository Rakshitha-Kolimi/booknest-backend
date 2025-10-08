package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"booknest/internal/domain"
)

type BookController struct {
	s domain.BookService
}

func NewBookController(s domain.BookService) BookController {
	return BookController{
		s: s,
	}
}

func GetHealth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (c BookController) GetBooks(ctx *gin.Context) {
	books, err := c.s.GetBooks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the data
	ctx.JSON(http.StatusOK, books)
}

func (c BookController) GetBook(ctx *gin.Context) {
	// Get the ID
	id := ctx.Param("id")

	// Parse the ID
	bookId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := c.s.GetBook(ctx, bookId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the data
	ctx.JSON(http.StatusOK, book)
}

func (c BookController) AddBook(ctx *gin.Context) {
	// Bind the input
	var in domain.BookInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := c.s.CreateBook(ctx, in)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the data
	ctx.JSON(http.StatusOK, book)
}

func (c BookController) UpdateBook(ctx *gin.Context) {
	// Get the ID
	id := ctx.Param("id")

	// Parse the ID
	bookId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Bind the input
	var in domain.BookInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := c.s.UpdateBook(ctx, bookId, in)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the data
	ctx.JSON(http.StatusOK, book)
}

func (c BookController) DeleteBook(ctx *gin.Context) {
	// Get the ID
	id := ctx.Param("id")

	// Parse the ID
	bookId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.s.DeleteBook(ctx, bookId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the data
	successMsg := fmt.Sprintf("Book with id %s is successfully deleted!", id)
	ctx.JSON(http.StatusOK, successMsg)
}
