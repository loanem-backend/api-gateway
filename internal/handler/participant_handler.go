package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/internal/dto"
	"github.com/loanem-backend/api-gateway/pkg/respx"
	pbparticipant "github.com/loanem-backend/protos/pb/proto/services/participant/v1"
)

type ParticipantHandler struct {
	teamClient        pbparticipant.TeamServiceClient
	participantClient pbparticipant.ParticipantServiceClient
}

func NewParticipantHandler(tc pbparticipant.TeamServiceClient, pc pbparticipant.ParticipantServiceClient) *ParticipantHandler {
	return &ParticipantHandler{
		teamClient:        tc,
		participantClient: pc,
	}
}

func (h *ParticipantHandler) CreateClasses(c *gin.Context) {
	idParam, err := parseIntParam(c, "courseId")
	if err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid param", err))
		return
	}

	var req pbparticipant.AddClassesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("failed parsing body", err))
		return
	}

	req.CourseId = idParam

	ctx := setLoginDataToContext(c)

	if _, err := h.teamClient.AddClasses(ctx, &req); err != nil {
		if c.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, respx.ResponseFail("service timeout", c.Err()))
			return
		}

		c.JSON(http.StatusInternalServerError, respx.ResponseFail("failed adding classes", err))
		return
	}

	c.JSON(http.StatusCreated, respx.ResponseSucceed("Classes successfully created", nil))
}

func (h *ParticipantHandler) GetClassesByCourse(c *gin.Context) {
	idParam, err := parseIntParam(c, "courseId")
	if err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid param", err))
		return
	}

	req := pbparticipant.GetClassesByCourseIDRequest{
		CourseId: idParam,
	}

	resp, err := h.teamClient.GetClassesByCourseID(c, &req)
	if err != nil {
		if c.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, respx.ResponseFail("service timeout", c.Err()))
			return
		}

		c.JSON(http.StatusInternalServerError, respx.ResponseFail("failed fetching classes", err))
		return
	}

	c.JSON(http.StatusOK, respx.ResponseSucceed("Classes successfully retrieved", dto.GetClassesByCourseIDResponseToClassResponses(resp)))
}
