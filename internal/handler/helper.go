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

func parseIntParam(c *gin.Context, param string) (int, error) {
	value := c.Param(param)

	return strconv.Atoi(value)
}
