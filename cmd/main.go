package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/loanem-backend/api-gateway/internal/handler"
)

func main() {
	godotenv.Load()

	r := gin.Default()

	handler.Start(r)
}
