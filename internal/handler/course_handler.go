package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/internal/dto"
	"github.com/loanem-backend/api-gateway/pkg/respx"
	pbcourse "github.com/loanem-backend/protos/pb/proto/services/course/v1"
)

type CourseHandler struct {
	courseClient pbcourse.CourseServiceClient
}

func NewCourseHandler(cc pbcourse.CourseServiceClient) *CourseHandler {
	return &CourseHandler{
		courseClient: cc,
	}
}

func (h *CourseHandler) Create(c *gin.Context) {
	var req dto.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail(messageInvalidBody, err))
		return
	}

	ctx := setLoginDataToContext(c)

	resp, err := h.courseClient.AddCourse(ctx, dto.CreateCourseRequestDTOToPB(&req))
	if err != nil {
		if err == context.DeadlineExceeded || c.Request.Context().Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, respx.ResponseFail(messageServiceTimeout, err))
			return
		}

		c.JSON(http.StatusInternalServerError, respx.ResponseFail(messageCreateCourseFailed, err))
		return
	}

	c.JSON(http.StatusCreated, respx.ResponseSucceed(messageCreateCourseSucceed, dto.NewCreateCourseResponse(resp)))
}

func (h *CourseHandler) Remove(c *gin.Context) {
	idParam, err := parseIntParam(c, "courseId")
	if err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid param", err))
		return
	}

	req := pbcourse.RemoveCourseRequest{
		Id: idParam,
	}

	ctx := setLoginDataToContext(c)

	if _, err := h.courseClient.RemoveCourse(ctx, &req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("failed deleting course", err))
		return
	}

	c.JSON(http.StatusOK, respx.ResponseSucceed("Course successfully deleted", nil))
}
