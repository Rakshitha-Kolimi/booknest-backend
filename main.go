package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"booknest/internal/http/controller"
	"booknest/internal/http/routes"
	"booknest/internal/middleware"
)

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.ErrorHandler())

	r.GET(routes.HealthRoute, controller.GetHealth)
	r.GET(routes.BooksRoute, controller.GetBooks)
	r.POST(routes.BookRoute, controller.AddBook)
	r.GET(routes.UsersRoute, controller.GetUsers)
	r.POST(routes.UserRoute, controller.AddUser)
	
	http.ListenAndServe(":8080", r)
}
