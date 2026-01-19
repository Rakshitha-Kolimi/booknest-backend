package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"booknest/internal/domain"
)

type bookController struct {
	service domain.BookService
}

func NewBookController(service domain.BookService) domain.BookController {
	return &bookController{service: service}
}

func (c *bookController) RegisterRoutes(r *gin.Engine) {
	books := r.Group("/books")
	{
		books.POST("", c.createBook)
		books.POST("/filter", c.filterBooks)
		books.GET("/:id", c.getBook)
		books.GET("", c.listBooks)
	}
}

func (c *bookController) createBook(ctx *gin.Context) {
	var input domain.BookInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := c.service.CreateBook(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, book)
}

func (c *bookController) getBook(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	book, err := c.service.GetBook(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

func (c *bookController) listBooks(ctx *gin.Context) {
	limit := 10
	offset := 0

	books, err := c.service.ListBooks(ctx, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (c *bookController) filterBooks(ctx *gin.Context) {
	var filter domain.BookFilter

	if v := ctx.Query("search"); v != "" {
		filter.Search = &v
	}

	if err := ctx.ShouldBindJSON(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := uint64(10)
	offset := uint64(0)

	if v := ctx.Query("limit"); v != "" {
		limit, _ = strconv.ParseUint(v, 10, 64)
	}

	if v := ctx.Query("offset"); v != "" {
		offset, _ = strconv.ParseUint(v, 10, 64)
	}

	result, err := c.service.FilterByCriteria(
		ctx,
		filter,
		domain.QueryOptions{Limit: limit, Offset: offset},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
