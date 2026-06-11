package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/pkg/respx"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
)

func Auth(ac pbauth.AuthServiceClient, authFn func(c *gin.Context) (string, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := authFn(ctx)
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
