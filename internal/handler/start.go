package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/internal/middleware"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
	"google.golang.org/grpc"
)

func Start(ge *gin.Engine, ac *grpc.ClientConn) {
	var (
		authClient      = pbauth.NewAuthServiceClient(ac)
		assistantClient = pbauth.NewAssistantServiceClient(ac)
	)

	ah := initHandlers(authClient, assistantClient)

	r := ge.Use(middleware.Timeout(2 * time.Second))
	routes(r, authClient, ah)
}

func initHandlers(ac pbauth.AuthServiceClient, asc pbauth.AssistantServiceClient) *AuthHandler {
	authHand := NewAuthHandler(ac, asc)

	return authHand
}

func routes(r gin.IRoutes, ac pbauth.AuthServiceClient, ah *AuthHandler) {
	r.POST("/login", ah.Login)
	r.POST("/assistants", ah.Create)
	r.PATCH("/me/password", middleware.Auth(ac), ah.SetPassword)
}
