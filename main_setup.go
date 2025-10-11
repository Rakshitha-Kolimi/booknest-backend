package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"booknest/internal/http/controller"
	"booknest/internal/http/database"
	"booknest/internal/http/routes"
	"booknest/internal/middleware"
	"booknest/internal/repository"
	"booknest/internal/service"
)

func setupServer() (*gin.Engine, error) {
	dbpool, err := database.Connect()
	if err != nil {
		return nil, err
	}

	bookRepo := repository.NewBookRepositoryImpl(dbpool)
	bookService := service.NewBookServiceImpl(bookRepo)
	bookController := controller.NewBookController(bookService)

	userRepo := repository.NewUserRepositoryImpl(dbpool)
	userService := service.NewUserServiceImpl(userRepo)
	userController := controller.NewUserController(userService)

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
		auth.DELETE(routes.BookIDRoute, bookController.DeleteBook)

		auth.GET(routes.UsersRoute, userController.GetUsers)
		auth.GET(routes.UserRoute, userController.GetUserByID)
		auth.DELETE(routes.UserRoute, userController.DeleteUser)
	}

	r.POST(routes.RegisterRoute, userController.RegisterUser)
	r.POST(routes.LoginRoute, userController.LoginUser)

	return r, nil
}

// StartHTTPServer starts the HTTP server — only used by main.go
func startHTTPServer(r *gin.Engine) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	log.Println("🚀 BookNest backend started on http://localhost:8080")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
