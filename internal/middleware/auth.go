package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/pkg/respx"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
)

func Auth(ac pbauth.AuthServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := GetAuthorization(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, respx.ResponseFail(
				"unauthorized", err,
			))
			return
		}

		resp, err := ac.ValidateToken(
			ctx,
			&pbauth.ValidateTokenRequest{
				AccessToken: accessToken,
			},
		)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, respx.ResponseFail(
				"unauthorized",
				errors.New("invalid or expired token"),
			))
			return
		}

		ctx.Set("id", resp.GetUserId())
		ctx.Set("name", resp.GetName())

		ctx.Next()
	}
}

func GetAuthorization(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("missing authorization header")
	}

	parts := strings.SplitN(header, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("invalid token")
	}

	return parts[1], nil
}
