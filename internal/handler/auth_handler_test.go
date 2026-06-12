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
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
