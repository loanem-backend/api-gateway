package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	server_mock "github.com/loanem-backend/api-gateway/internal/mocks/server"
	pbcourse "github.com/loanem-backend/protos/pb/proto/services/course/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCourseHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		mockBehavior func(m *server_mock.MockCourseServiceClient)
		body         any
		code         int
	}{
		{
			name: "Success - Created",
			mockBehavior: func(m *server_mock.MockCourseServiceClient) {
				m.EXPECT().
					AddCourse(gomock.Any(), gomock.Any()).
					Return(&pbcourse.AddCourseResponse{
						Id: 51,
					}, nil)
			},
			body: &pbcourse.AddCourseRequest{
				Name: "Course Test",
				Year: 1980,
			},
			code: 201,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCourseClient := server_mock.NewMockCourseServiceClient(ctrl)
			test.mockBehavior(mockCourseClient)

			h := NewCourseHandler(mockCourseClient)

			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			r.Use(mockAuthMiddleware())

			r.POST("/courses", h.Create)

			var b bytes.Buffer
			if strBody, ok := test.body.(string); ok {
				b.WriteString(strBody)
			} else {
				_ = json.NewEncoder(&b).Encode(test.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/courses", &b)
			req.Header.Set(strContentType, strApplicationJSON)
			c.Request = req

			r.ServeHTTP(w, req)

			assert.Equal(t, test.code, w.Code)
		})
	}
}
