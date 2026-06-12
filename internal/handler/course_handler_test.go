package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/internal/dto"
	server_mock "github.com/loanem-backend/api-gateway/internal/mocks/server"
	pbcourse "github.com/loanem-backend/protos/pb/proto/services/course/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCourseHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validBody := &dto.CreateCourseRequest{
		Name: "Course Test",
		Year: 2025,
	}

	tests := []struct {
		name         string
		mockBehavior func(m *server_mock.MockCourseServiceClient)
		body         any
		mutateReq    func(req *http.Request) *http.Request
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
			body: validBody,
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, w.Code)
				assert.Contains(t, w.Body.String(), messageCreateCourseSucceed)
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
				assert.Contains(t, w.Body.String(), messageInvalidBody)
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
				assert.Contains(t, w.Body.String(), messageInvalidBody)
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
				assert.Contains(t, w.Body.String(), messageInvalidBody)
			},
		},
		{
			name: "Failed_RequestTimeout",
			mockBehavior: func(m *server_mock.MockCourseServiceClient) {
				m.EXPECT().AddCourse(gomock.Any(), gomock.Any()).Return(nil, errors.New(""))
			},
			body: validBody,
			mutateReq: func(req *http.Request) *http.Request {
				ctx, cancel := context.WithTimeout(req.Context(), 0)
				defer cancel()

				return req.WithContext(ctx)
			},
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusRequestTimeout, w.Code)
				strBody := w.Body.String()
				assert.NotContains(t, strBody, messageInvalidBody)
				assert.Contains(t, strBody, messageRequestTimeout)
			},
		},
		{
			name: "Failed_RequestCanceled",
			mockBehavior: func(m *server_mock.MockCourseServiceClient) {
				m.EXPECT().AddCourse(gomock.Any(), gomock.Any()).Return(nil, context.Canceled)
			},
			body: validBody,
			mutateReq: func(req *http.Request) *http.Request {
				ctx, cancel := context.WithCancel(req.Context())
				cancel()

				return req.WithContext(ctx)
			},
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, httpStatusClientClosedRequest, w.Code)
				strBody := w.Body.String()
				assert.NotContains(t, strBody, messageInvalidBody)
				assert.NotContains(t, strBody, messageRequestTimeout)
				assert.Contains(t, strBody, messageClientClosedRequest)
			},
		},
		{
			name: "Failed_gRPCTimeout",
			mockBehavior: func(m *server_mock.MockCourseServiceClient) {
				m.EXPECT().
					AddCourse(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.DeadlineExceeded, messageServiceTimeout))
			},
			body: validBody,
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusGatewayTimeout, w.Code)
				strBody := w.Body.String()
				assert.NotContains(t, strBody, messageInvalidBody)
				assert.NotContains(t, strBody, messageRequestTimeout)
				assert.NotContains(t, strBody, messageClientClosedRequest)
				assert.Contains(t, strBody, messageServiceTimeout)
			},
		},
		{
			name: "Failed_InternalServerError",
			mockBehavior: func(m *server_mock.MockCourseServiceClient) {
				m.EXPECT().
					AddCourse(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, ""))
			},
			body: validBody,
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				strBody := w.Body.String()
				assert.NotContains(t, strBody, messageInvalidBody)
				assert.NotContains(t, strBody, messageServiceTimeout)
				assert.Contains(t, strBody, messageInternalServerError)
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

			if test.mutateReq != nil {
				req = test.mutateReq(req)
			}
			c.Request = req

			r.ServeHTTP(w, req)

			test.assertCase(t, w)
		})
	}
}

func BenchmarkCourseHandler_Create(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockCourseClient := server_mock.NewMockCourseServiceClient(ctrl)
	h := NewCourseHandler(mockCourseClient)

	mockCourseClient.EXPECT().
		AddCourse(gomock.Any(), gomock.Any()).
		Return(&pbcourse.AddCourseResponse{Id: 321}, nil).
		AnyTimes()

	validBody := []byte(`{"name":"Course Test","year":2025}`)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(mockAuthMiddleware())
	r.POST("/courses", h.Create)

	for range b.N {
		b.StopTimer()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(validBody))
		req.Header.Set(strContentType, strApplicationJSON)

		b.StartTimer()

		r.ServeHTTP(w, req)
	}
}
