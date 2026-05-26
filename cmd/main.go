package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/loanem-backend/api-gateway/internal/config"
	"github.com/loanem-backend/api-gateway/internal/handler"
)

func main() {
	_ = godotenv.Load()

	r := gin.Default()

	authConn, courseConn, inventoryConn := handler.InitConnections()

	handler.Start(r, authConn, courseConn, inventoryConn)

	port := config.GetEnv("APP_PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Printf("Api gateway running on port %s\n", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("forced to shutdown: ", err)
	}

	handler.CloseConnections(authConn, courseConn, inventoryConn)

	log.Println("Exited")
}
