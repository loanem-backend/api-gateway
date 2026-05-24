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
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, respx.ResponseFail(
				"unauthorized",
				errors.New("missing authorization header"),
			))
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, respx.ResponseFail(
				"unauthorized",
				errors.New("invalid token"),
			))
			return
		}

		resp, err := ac.ValidateToken(
			ctx,
			&pbauth.ValidateTokenRequest{
				Token: parts[1],
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
