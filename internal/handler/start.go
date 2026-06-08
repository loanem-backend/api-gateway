package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/internal/middleware"
	"github.com/loanem-backend/api-gateway/internal/service"
	"github.com/loanem-backend/api-gateway/pkg/storage"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
	pbcourse "github.com/loanem-backend/protos/pb/proto/services/course/v1"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	pbparticipant "github.com/loanem-backend/protos/pb/proto/services/participant/v1"
	"google.golang.org/grpc"
)

func Start(ge *gin.Engine, ac *grpc.ClientConn, cc *grpc.ClientConn, ic *grpc.ClientConn, pc *grpc.ClientConn, sc *storage.S3Client) {
	var (
		authClient       = pbauth.NewAuthServiceClient(ac)
		assistantClient  = pbauth.NewAssistantServiceClient(ac)
		courseClient     = pbcourse.NewCourseServiceClient(cc)
		instrumentClient = pbinventory.NewInstrumentServiceClient(ic)
		toolkitClient    = pbinventory.NewToolkitServiceClient(ic)
		teamClient       = pbparticipant.NewTeamServiceClient(pc)
	)

	ah, ch, ih, ph := initHandlers(authClient, assistantClient, courseClient, instrumentClient, toolkitClient, teamClient, sc)

	r := ge.Use(middleware.Timeout(2 * time.Second))
	routes(r, authClient, ah, ch, ih, ph)
}

func initHandlers(
	ac pbauth.AuthServiceClient,
	asc pbauth.AssistantServiceClient,
	cc pbcourse.CourseServiceClient,
	ic pbinventory.InstrumentServiceClient,
	tc pbinventory.ToolkitServiceClient,
	tmc pbparticipant.TeamServiceClient,
	strc *storage.S3Client,
) (*AuthHandler, *CourseHandler, *InventoryHandler, *ParticipantHandler) {
	var (
		storageServ = service.NewStorageService(strc)
	)

	var (
		authHand        = NewAuthHandler(ac, asc)
		courseHand      = NewCourseHandler(cc)
		inventoryHand   = NewInventoryHandler(ic, tc, storageServ)
		participantHand = NewParticipantHandler(tmc, nil)
	)

	return authHand, courseHand, inventoryHand, participantHand
}

func routes(r gin.IRoutes, ac pbauth.AuthServiceClient, ah *AuthHandler, ch *CourseHandler, ih *InventoryHandler, ph *ParticipantHandler) {
	r.POST("/login", ah.Login)
	r.POST("/assistants", ah.Create)
	r.PATCH("/me/password", middleware.Auth(ac), ah.SetPassword)

	r.POST("/courses", middleware.Auth(ac), ch.Create)
	r.DELETE("/courses/:courseId", middleware.Auth(ac), ch.Remove)

	r.POST("/toolkits", middleware.Auth(ac), ih.CreateToolkit)
	r.POST("/instruments", middleware.Auth(ac), ih.CreateInstrument)
	r.PATCH("/instruments/:instrumentId/picture", middleware.Auth(ac), ih.SetInstrumentPicture)

	r.POST("/courses/:courseId/classes", middleware.Auth(ac), ph.CreateClasses)
	r.GET("/courses/:courseId/classes", ph.GetClassesByCourse)
}
