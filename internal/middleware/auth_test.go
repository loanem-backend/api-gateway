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
)

func TestAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validToken := "sample.valid.token"
	fullToken := "Bearer " + validToken

	tests := []struct {
		name         string
		mockBehavior func(m *server_mock.MockAuthServiceClient)
		assertCase   func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "Success_ValidToken",
			mockBehavior: func(m *server_mock.MockAuthServiceClient) {
				m.EXPECT().
					ValidateToken(gomock.Any(), &pbauth.ValidateTokenRequest{
						AccessToken: validToken,
					}).
					Return(&pbauth.ValidateTokenResponse{
						UserId:  21,
						Name:    "Name Test",
						IsValid: true,
					}, nil)
			},
			assertCase: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)
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
			r.GET("/sample-path", Auth(mockAuthClient), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"ok": true,
				})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/sample-path", nil)
			req.Header.Set("Authorization", fullToken)

			r.ServeHTTP(w, req)

			test.assertCase(t, w)
		})
	}
}
