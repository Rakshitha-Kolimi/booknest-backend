package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

type DummyBook struct {
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
}

var dummyData = []DummyBook{
	{
		Title:  "The Pragmatic Programmer",
		Author: "Andrew Hunt, David Thomas",
		Price:  42.50,
		Stock:  12,
	}, {
		Title:  "Clean Code",
		Author: "Robert C. Martin",
		Price:  39.99,
		Stock:  8,
	}, {
		Title:  "Introduction to Algorithms",
		Author: "Thomas H. Cormen",
		Price:  89.95,
		Stock:  5,
	}, {
		Title:  "Design Patterns: Elements of Reusable Object-Oriented Software",
		Author: "Erich Gamma, Richard Helm, Ralph Johnson, John Vlissides",
		Price:  55.00,
		Stock:  7,
	}, {
		Title:  "You Don’t Know JS Yet",
		Author: "Kyle Simpson",
		Price:  25.00,
		Stock:  20,
	}}

func GetHealth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GetBooks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dummyData)
}

func (c BookController) AddBook(ctx *gin.Context) {
	// Bind the input
	var in domain.BookInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book,err := c.s.CreateBook(in)
	if  err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the data
	ctx.JSON(http.StatusOK, book)
}
