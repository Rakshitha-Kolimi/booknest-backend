package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"booknest/internal/http/controller"
	"booknest/internal/http/database"
	"booknest/internal/http/routes"
	"booknest/internal/middleware"
	"booknest/internal/repository"
	"booknest/internal/service"
)

func main() {
	godotenv.Load()
	dbpool, err:= database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	bookRepo := repository.NewBookRepositoryImpl(dbpool)
	bookService := service.NewBookServiceImpl(bookRepo)
	bookController := controller.NewBookController(bookService)

	log.Println("BookNest backend started...")

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.ErrorHandler())

	r.GET(routes.HealthRoute, controller.GetHealth)
	r.GET(routes.BooksRoute, controller.GetBooks)
	r.POST(routes.BookRoute, bookController.AddBook)

	r.GET(routes.UsersRoute, controller.GetUsers)
	r.POST(routes.UserRoute, controller.AddUser)

	http.ListenAndServe(":8080", r)
}
