package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func setLoginDataToContext(c *gin.Context) context.Context {
	md := metadata.Pairs(
		"id", strconv.Itoa(int(c.MustGet("id").(int32))),
		"name", c.MustGet("name").(string),
	)

	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	return ctx
}

func parseIntParam(c *gin.Context, param string) (int32, error) {
	value := c.Param(param)

	intVal, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(intVal), nil
}

const (
	strContentType     = "Content-Type"
	strApplicationJSON = "application/json"
)

func mockAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("id", int32(99))
		ctx.Set("name", "Arvin Fernanda")
		ctx.Next()
	}
}
