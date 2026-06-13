package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	server_mock "github.com/loanem-backend/api-gateway/internal/mocks/server"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validToken := "sample.valid.token"
	fullToken := "Bearer " + validToken

	sampleResponse := &pbauth.ValidateTokenResponse{
		UserId:  21,
		Name:    "Name Test",
		IsValid: true,
	}

	tests := []struct {
		name         string
		mockBehavior func(m *server_mock.MockAuthServiceClient)
		headerValue  string
		assertCase   func(t *testing.T, w *httptest.ResponseRecorder, c *gin.Context)
	}{
		{
			name: "Success_ValidToken",
			mockBehavior: func(m *server_mock.MockAuthServiceClient) {
				m.EXPECT().
					ValidateToken(gomock.Any(), &pbauth.ValidateTokenRequest{
						AccessToken: validToken,
					}).
					Return(sampleResponse, nil)
			},
			headerValue: fullToken,
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder, c *gin.Context) {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.NotContains(t, w.Body.String(), "unauthorized")
				id, exists := c.Get("id")
				assert.True(t, exists)
				assert.Equal(t, sampleResponse.UserId, id)
				name, exists := c.Get("name")
				assert.True(t, exists)
				assert.Equal(t, sampleResponse.Name, name)
			},
		},
		{
			name: "Failed_MissingHeader",
			mockBehavior: func(m *server_mock.MockAuthServiceClient) {
				m.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).Times(0)
			},
			headerValue: "",
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder, c *gin.Context) {
				assert.Equal(t, http.StatusUnauthorized, w.Code)
				strBody := w.Body.String()
				assert.Contains(t, strBody, "unauthorized")
				assert.NotContains(t, strBody, "invalid or expired token")
			},
		},
		{
			name: "Failed_InvalidToken",
			mockBehavior: func(m *server_mock.MockAuthServiceClient) {
				m.EXPECT().
					ValidateToken(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, ""))
			},
			headerValue: fullToken,
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder, c *gin.Context) {
				assert.Equal(t, http.StatusUnauthorized, w.Code)
				assert.Contains(t, w.Body.String(), "invalid or expired token")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthClient := server_mock.NewMockAuthServiceClient(ctrl)
			test.mockBehavior(mockAuthClient)

			r := gin.New()
			var ctx *gin.Context
			r.GET("/sample-path", Auth(mockAuthClient), func(c *gin.Context) {
				ctx = c

				c.JSON(http.StatusOK, gin.H{
					"ok": true,
				})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/sample-path", nil)
			req.Header.Set("Authorization", test.headerValue)

			r.ServeHTTP(w, req)

			test.assertCase(t, w, ctx)
		})
	}
}

func TestGetAuthorization(t *testing.T) {
	sampleToken := "sample.token"

	tests := []struct {
		name        string
		headerValue string
		assertCase  func(*testing.T, string, error)
	}{
		{
			name:        "Success",
			headerValue: "Bearer " + sampleToken,
			assertCase: func(t *testing.T, s string, err error) {
				assert.Equal(t, sampleToken, s)
				assert.NoError(t, err)
			},
		},
		{
			name:        "Failed_MissingHeader",
			headerValue: "",
			assertCase: func(t *testing.T, s string, err error) {
				assert.Equal(t, "", s)
				assert.Error(t, err)
			},
		},
		{
			name:        "Failed_InvalidScheme",
			headerValue: "My " + sampleToken,
			assertCase: func(t *testing.T, s string, err error) {
				assert.NotEqual(t, sampleToken, s)
				assert.Equal(t, "", s)
				assert.Error(t, err)
			},
		},
		{
			name:        "Failed_Malformed",
			headerValue: "Bearer" + sampleToken,
			assertCase: func(t *testing.T, s string, err error) {
				assert.NotEqual(t, sampleToken, s)
				assert.Equal(t, "", s)
				assert.Error(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.headerValue != "" {
				req.Header.Set("Authorization", test.headerValue)
			}
			c.Request = req

			token, err := GetAuthorization(c)

			test.assertCase(t, token, err)
		})
	}
}
