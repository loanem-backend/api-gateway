package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/internal/middleware"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"google.golang.org/grpc"
)

func Start(ge *gin.Engine, ac *grpc.ClientConn, ic *grpc.ClientConn) {
	var (
		authClient      = pbauth.NewAuthServiceClient(ac)
		assistantClient = pbauth.NewAssistantServiceClient(ac)
		inventoryClient = pbinventory.NewInventoryServiceClient(ic)
	)

	ah, ih := initHandlers(authClient, assistantClient, inventoryClient)

	r := ge.Use(middleware.Timeout(2 * time.Second))
	routes(r, authClient, ah, ih)
}

func initHandlers(
	ac pbauth.AuthServiceClient,
	asc pbauth.AssistantServiceClient,
	ic pbinventory.InventoryServiceClient,
) (*AuthHandler, *InventoryHandler) {
	authHand := NewAuthHandler(ac, asc)
	inventoryHand := NewInventoryHandler(ic)

	return authHand, inventoryHand
}

func routes(r gin.IRoutes, ac pbauth.AuthServiceClient, ah *AuthHandler, ih *InventoryHandler) {
	r.POST("/login", ah.Login)
	r.POST("/assistants", ah.Create)
	r.PATCH("/me/password", middleware.Auth(ac), ah.SetPassword)

	r.POST("/instruments", middleware.Auth(ac), ih.Create)
}
