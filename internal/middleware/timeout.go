package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(ctx, d)
		defer cancel()

		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}
