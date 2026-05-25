package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/internal/middleware"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
	pbcourse "github.com/loanem-backend/protos/pb/proto/services/course/v1"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"google.golang.org/grpc"
)

func Start(ge *gin.Engine, ac *grpc.ClientConn, cc *grpc.ClientConn, ic *grpc.ClientConn) {
	var (
		authClient      = pbauth.NewAuthServiceClient(ac)
		assistantClient = pbauth.NewAssistantServiceClient(ac)
		courseClient    = pbcourse.NewCourseServiceClient(cc)
		inventoryClient = pbinventory.NewInventoryServiceClient(ic)
	)

	ah, ch, ih := initHandlers(authClient, assistantClient, courseClient, inventoryClient)

	r := ge.Use(middleware.Timeout(2 * time.Second))
	routes(r, authClient, ah, ch, ih)
}

func initHandlers(
	ac pbauth.AuthServiceClient,
	asc pbauth.AssistantServiceClient,
	cc pbcourse.CourseServiceClient,
	ic pbinventory.InventoryServiceClient,
) (*AuthHandler, *CourseHandler, *InventoryHandler) {
	authHand := NewAuthHandler(ac, asc)
	courseHand := NewCourseHandler(cc)
	inventoryHand := NewInventoryHandler(ic)

	return authHand, courseHand, inventoryHand
}

func routes(r gin.IRoutes, ac pbauth.AuthServiceClient, ah *AuthHandler, ch *CourseHandler, ih *InventoryHandler) {
	r.POST("/login", ah.Login)
	r.POST("/assistants", ah.Create)
	r.PATCH("/me/password", middleware.Auth(ac), ah.SetPassword)

	r.POST("/courses", middleware.Auth(ac), ch.Create)

	r.POST("/instruments", middleware.Auth(ac), ih.Create)
}
