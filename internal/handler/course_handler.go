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
	var req pbcourse.AddCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid body", err))
		return
	}

	ctx := setLoginDataToContext(c)

	resp, err := h.courseClient.AddCourse(ctx, &req)
	if err != nil {
		if c.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, respx.ResponseFail("service timeout", c.Err()))
			return
		}

		c.JSON(http.StatusInternalServerError, respx.ResponseFail("failed creating course", err))
		return
	}

	c.JSON(http.StatusCreated, respx.ResponseSucceed("Course successfully created", dto.NewCreateCourseResponse(resp)))
}

func (h *CourseHandler) Remove(c *gin.Context) {
	idParam, err := parseIntParam(c, "courseId")
	if err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid param", err))
		return
	}

	req := pbcourse.RemoveCourseRequest{
		Id: int32(idParam),
	}

	ctx := setLoginDataToContext(c)

	if _, err := h.courseClient.RemoveCourse(ctx, &req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("failed deleting course", err))
		return
	}

	c.JSON(http.StatusOK, respx.ResponseSucceed("Course successfully deleted", nil))
}
