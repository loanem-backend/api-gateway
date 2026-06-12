package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/pkg/respx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func handleErrorFromClient(c *gin.Context, err error) {
	cErr := c.Request.Context().Err()
	if errors.Is(cErr, context.DeadlineExceeded) {
		c.JSON(http.StatusRequestTimeout, respx.ResponseFail(messageRequestTimeout, cErr))
		return
	}
	if errors.Is(cErr, context.Canceled) {
		c.JSON(httpStatusClientClosedRequest, respx.ResponseFail(messageClientClosedRequest, cErr))
		return
	}

	if errors.Is(err, context.DeadlineExceeded) || status.Code(err) == codes.DeadlineExceeded {
		c.JSON(http.StatusGatewayTimeout, respx.ResponseFail(messageServiceTimeout, err))
		return
	}

	c.JSON(http.StatusInternalServerError, respx.ResponseFail(messageInternalServerError, err))
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
