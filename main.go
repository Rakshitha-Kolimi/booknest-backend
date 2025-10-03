package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"booknest/internal/http/controller"
	"booknest/internal/http/routes"
)

func main() {
	r := gin.Default()
	r.GET(routes.HealthRoute, controller.GetHealth)
	r.GET(routes.BooksRoute, controller.GetBooks)
	r.POST(routes.BookRoute, controller.AddBook)

	http.ListenAndServe(":8080", r)
}
