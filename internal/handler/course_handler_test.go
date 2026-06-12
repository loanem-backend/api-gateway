package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/internal/dto"
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
		assertCase   func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "Success_Created",
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
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, w.Code)
			},
		},
		{
			name: "Failed_InvalidBody",
			mockBehavior: func(m *server_mock.MockCourseServiceClient) {
				m.EXPECT().
					AddCourse(gomock.Any(), gomock.Any()).
					Times(0)
			},
			body: "raw string",
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Contains(t, w.Body.String(), "invalid body")
			},
		},
		{
			name: "Failed_MissingNameField",
			mockBehavior: func(m *server_mock.MockCourseServiceClient) {
				m.EXPECT().
					AddCourse(gomock.Any(), gomock.Any()).
					Times(0)
			},
			body: &dto.CreateCourseRequest{
				Year: 2025,
			},
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Contains(t, w.Body.String(), "invalid body")
			},
		},
		{
			name: "Failed_MissingYearField",
			mockBehavior: func(m *server_mock.MockCourseServiceClient) {
				m.EXPECT().
					AddCourse(gomock.Any(), gomock.Any()).
					Times(0)
			},
			body: &dto.CreateCourseRequest{
				Name: "Course Test",
			},
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Contains(t, w.Body.String(), "invalid body")
			},
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

			test.assertCase(t, w)
		})
	}
}
