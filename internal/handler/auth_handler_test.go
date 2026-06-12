package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/internal/dto"
	server_mock "github.com/loanem-backend/api-gateway/internal/mocks/server"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validBody := &dto.LoginRequest{
		Phone:    "089089089089",
		Password: "drowssap",
	}

	clientOkResponse := &pbauth.LoginResponse{
		AccessToken:           "abcd.efghijkl.mnop",
		RefreshToken:          "qwerty.asdfghjkl.zxcv",
		RefreshExpirationHour: 12,
	}

	tests := []struct {
		name         string
		mockBehavior func(m *server_mock.MockAuthServiceClient)
		body         any
		mutateReq    func(req *http.Request) *http.Request
		assertCase   func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "Success_Created",
			mockBehavior: func(m *server_mock.MockAuthServiceClient) {
				m.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Return(clientOkResponse, nil)
			},
			body: validBody,
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, w.Code)
				assert.Contains(t, w.Body.String(), messageLoginSucceed)
				assert.Contains(t, w.Header().Get("Set-Cookie"), "refresh_token=")
			},
		},
		{
			name: "Failed_InvalidBody",
			mockBehavior: func(m *server_mock.MockAuthServiceClient) {
				m.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Times(0)
			},
			body: "raw string",
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Contains(t, w.Body.String(), messageInvalidBody)
			},
		},
		{
			name: "Failed_MissingPasswordField",
			mockBehavior: func(m *server_mock.MockAuthServiceClient) {
				m.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Times(0)
			},
			body: &dto.LoginRequest{
				Phone: "089089089089",
			},
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Contains(t, w.Body.String(), messageInvalidBody)
			},
		},
		{
			name: "Failed_MissingPhoneField",
			mockBehavior: func(m *server_mock.MockAuthServiceClient) {
				m.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Times(0)
			},
			body: &dto.LoginRequest{
				Password: "pawwsord",
			},
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Contains(t, w.Body.String(), messageInvalidBody)
			},
		},
		{
			name: "Failed_RequestTimeout",
			mockBehavior: func(m *server_mock.MockAuthServiceClient) {
				m.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req any, opts ...any) (*pbauth.LoginResponse, error) {
						<-ctx.Done()
						return nil, ctx.Err()
					})
			},
			body: validBody,
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusRequestTimeout, w.Code)
				strBody := w.Body.String()
				assert.NotContains(t, strBody, messageInvalidBody)
				assert.Contains(t, strBody, messageRequestTimeout)
			},
		},
		{
			name: "Failed_RequestCanceled",
			mockBehavior: func(m *server_mock.MockAuthServiceClient) {
				m.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Return(nil, context.Canceled)
			},
			body: validBody,
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, httpStatusClientClosedRequest, w.Code)
				strBody := w.Body.String()
				assert.NotContains(t, strBody, messageInvalidBody)
				assert.NotContains(t, strBody, messageRequestTimeout)
				assert.Contains(t, strBody, messageClientClosedRequest)
			},
			mutateReq: func(req *http.Request) *http.Request {
				ctx, cancel := context.WithCancel(req.Context())
				cancel()

				return req.WithContext(ctx)
			},
		},
		{
			name: "Failed_gRPCTimeout",
			mockBehavior: func(m *server_mock.MockAuthServiceClient) {
				m.EXPECT().
					Login(gomock.Any(), gomock.Any()).
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
			mockBehavior: func(m *server_mock.MockAuthServiceClient) {
				m.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, ""))
			},
			body: validBody,
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				strBody := w.Body.String()
				assert.NotContains(t, strBody, messageInvalidBody)
				assert.NotContains(t, strBody, messageRequestTimeout)
				assert.NotContains(t, strBody, messageClientClosedRequest)
				assert.NotContains(t, strBody, messageServiceTimeout)
				assert.Contains(t, strBody, messageInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthClient := server_mock.NewMockAuthServiceClient(ctrl)
			test.mockBehavior(mockAuthClient)

			h := NewAuthHandler(mockAuthClient, nil)

			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			if strings.Contains(test.name, "RequestTimeout") {
				r.Use(func(c *gin.Context) {
					ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Millisecond)
					defer cancel()

					c.Request = c.Request.WithContext(ctx)
					c.Next()
				})
			}

			r.POST("/login", h.Login)

			var b bytes.Buffer
			if strBody, ok := test.body.(string); ok {
				b.WriteString(strBody)
			} else {
				_ = json.NewEncoder(&b).Encode(test.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/login", &b)
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
