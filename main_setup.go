package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"booknest/internal/http/controller"
	"booknest/internal/http/database"
	"booknest/internal/middleware"
	"booknest/internal/repository"
	"booknest/internal/service/author_service"
	"booknest/internal/service/book_service"
	"booknest/internal/service/cart_service"
	"booknest/internal/service/order_service"
	"booknest/internal/service/publisher_service"
	"booknest/internal/service/user_service"
)

func useCORSMiddleware(allowedOrigins map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)
			c.Writer.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, DELETE, OPTIONS",
			)
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func SetupServer(dbpool *pgxpool.Pool) (*gin.Engine, error) {
	gormdb, err := database.ConnectGORM()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := gormdb.DB()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepo(dbpool, gormdb)
	vtRepo := repository.NewVerificationRepo(dbpool, gormdb)
	userService := user_service.NewUserService(dbpool, userRepo, vtRepo)
	userController := controller.NewUserController(userService)

	bookRepo := repository.NewBookRepository(gormdb, sqlDB)
	bookService := book_service.NewBookService(bookRepo, gormdb)
	bookController := controller.NewBookController(bookService)

	authorRepo := repository.NewAuthorRepo(gormdb)
	authorService := author_service.NewAuthorService(authorRepo)
	authorController := controller.NewAuthorController(authorService)

	publisherRepo := repository.NewPublisherRepo(dbpool, gormdb)
	publisherService := publisher_service.NewPublisherService(dbpool, publisherRepo)
	publisherController := controller.NewPublisherController(publisherService)

	cartRepo := repository.NewCartRepo(dbpool)
	cartService := cart_service.NewCartService(dbpool, cartRepo, bookRepo)
	cartController := controller.NewCartController(cartService)

	orderRepo := repository.NewOrderRepo(dbpool)
	orderService := order_service.NewOrderService(dbpool, orderRepo, cartRepo)
	orderController := controller.NewOrderController(orderService)

	r := gin.Default()
	r.Use(useCORSMiddleware(map[string]bool{
		"http://localhost:3000": true,
		"http://localhost:5173": true,
	}))
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.ErrorHandler())
	r.GET(
		"/swagger/*any",
		middleware.SwaggerAuthMiddleware(),
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	userController.RegisterRoutes(r)
	bookController.RegisterRoutes(r)
	authorController.RegisterRoutes(r)
	publisherController.RegisterRoutes(r)
	cartController.RegisterRoutes(r)
	orderController.RegisterRoutes(r)

	return r, nil
}

// StartHTTPServer starts the HTTP server — only used by main.go
func StartHTTPServer(r *gin.Engine) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// If we don’t run it in a goroutine, shutdown logic will never execute
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	log.Println("🚀 BookNest backend started on http://localhost:8080")

	// graceful shutdown
	/*
		* Creates a channel to receive OS signals and Listens for:
			- Ctrl + C
			- Docker stop
			- Pod termination
			<-quit blocks until signal arrives
	*/
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	// Gives active requests 5 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	/*
		Why server shut down:
		1. Stops accepting new requests
		2. Waits for in-flight requests
		3. Closes idle connections
		4. Respects the timeout context
	*/
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
