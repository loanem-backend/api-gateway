package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/internal/dto"
	"github.com/loanem-backend/api-gateway/pkg/respx"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
)

type AuthHandler struct {
	client pbauth.AuthServiceClient
}

func NewAuthHandler(c pbauth.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		client: c,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var req pbauth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid body", err))
		return
	}

	resp, err := h.client.Login(ctx, &req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, respx.ResponseFail("service timeout", ctx.Err()))
			return
		}

		c.JSON(http.StatusInternalServerError, respx.ResponseFail("failed logging in", err))
		return
	}

	c.JSON(http.StatusCreated, respx.ResponseSucceed("Logged in successfully", dto.NewLoginResponse(resp)))
}
