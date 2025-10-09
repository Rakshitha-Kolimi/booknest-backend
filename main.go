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
	dbpool, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	bookRepo := repository.NewBookRepositoryImpl(dbpool)
	bookService := service.NewBookServiceImpl(bookRepo)
	bookController := controller.NewBookController(bookService)

	userRepo := repository.NewUserRepositoryImpl(dbpool)
	userService := service.NewUserServiceImpl(userRepo)
	userController := controller.NewUserController(userService)

	log.Println("BookNest backend started...")

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.ErrorHandler())

	r.GET(routes.HealthRoute, controller.GetHealth)

	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.POST(routes.BookRoute, bookController.AddBook)
		auth.GET(routes.BooksRoute, bookController.GetBooks)
		auth.GET(routes.BookIDRoute, bookController.GetBook)
		auth.PUT(routes.BookIDRoute, bookController.UpdateBook)
		auth.PUT(routes.BookIDRoute, bookController.DeleteBook)

		auth.GET(routes.UsersRoute, userController.GetUsers)
		auth.GET(routes.UserRoute, userController.GetUserByID)
		auth.DELETE(routes.UserRoute, userController.DeleteUser)
	}

	r.POST(routes.RegisterRoute, userController.RegisterUser)
	r.POST(routes.LoginRoute, userController.LoginUser)

	http.ListenAndServe(":8080", r)
}
