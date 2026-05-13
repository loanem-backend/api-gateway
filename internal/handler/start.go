package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Start(ge *gin.Engine) {
	authConn, err := grpc.NewClient(serviceAddr(auth), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("failed connecting to services: %w", err))
	}

	authClient := pbauth.NewAuthServiceClient(authConn)

	ah := initHandlers(authClient)

	routes(ge, ah)
}

func initHandlers(ac pbauth.AuthServiceClient) *AuthHandler {
	authHand := NewAuthHandler(ac)

	return authHand
}

func routes(r *gin.Engine, ah *AuthHandler) {
	r.POST("/login", ah.Login)
}
