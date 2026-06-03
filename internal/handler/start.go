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
		authClient       = pbauth.NewAuthServiceClient(ac)
		assistantClient  = pbauth.NewAssistantServiceClient(ac)
		courseClient     = pbcourse.NewCourseServiceClient(cc)
		instrumentClient = pbinventory.NewInstrumentServiceClient(ic)
		toolkitClient    = pbinventory.NewToolkitServiceClient(ic)
	)

	ah, ch, ih := initHandlers(authClient, assistantClient, courseClient, instrumentClient, toolkitClient)

	r := ge.Use(middleware.Timeout(2 * time.Second))
	routes(r, authClient, ah, ch, ih)
}

func initHandlers(
	ac pbauth.AuthServiceClient,
	asc pbauth.AssistantServiceClient,
	cc pbcourse.CourseServiceClient,
	ic pbinventory.InstrumentServiceClient,
	tc pbinventory.ToolkitServiceClient,
) (*AuthHandler, *CourseHandler, *InventoryHandler) {
	authHand := NewAuthHandler(ac, asc)
	courseHand := NewCourseHandler(cc)
	inventoryHand := NewInventoryHandler(ic, tc)

	return authHand, courseHand, inventoryHand
}

func routes(r gin.IRoutes, ac pbauth.AuthServiceClient, ah *AuthHandler, ch *CourseHandler, ih *InventoryHandler) {
	r.POST("/login", ah.Login)
	r.POST("/assistants", ah.Create)
	r.PATCH("/me/password", middleware.Auth(ac), ah.SetPassword)

	r.POST("/courses", middleware.Auth(ac), ch.Create)
	r.DELETE("/courses/:courseId", middleware.Auth(ac), ch.Remove)

	r.POST("/instruments", middleware.Auth(ac), ih.CreateInstrument)
	r.POST("/toolkits", middleware.Auth(ac), ih.CreateToolkit)
}
