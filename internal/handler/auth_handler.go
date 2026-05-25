package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/internal/dto"
	"github.com/loanem-backend/api-gateway/pkg/respx"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
)

type AuthHandler struct {
	authClient      pbauth.AuthServiceClient
	assistantClient pbauth.AssistantServiceClient
}

func NewAuthHandler(ac pbauth.AuthServiceClient, asc pbauth.AssistantServiceClient) *AuthHandler {
	return &AuthHandler{
		authClient:      ac,
		assistantClient: asc,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	// defer cancel()

	var req pbauth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid body", err))
		return
	}

	resp, err := h.authClient.Login(c, &req)
	if err != nil {
		if c.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, respx.ResponseFail("service timeout", c.Err()))
			return
		}

		c.JSON(http.StatusInternalServerError, respx.ResponseFail("failed logging in", err))
		return
	}

	c.JSON(http.StatusCreated, respx.ResponseSucceed("Logged in successfully", dto.NewLoginResponse(resp)))
}

func (h *AuthHandler) Create(c *gin.Context) {
	var req pbauth.CreateAssistantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid body", err))
		return
	}

	resp, err := h.assistantClient.CreateAssistant(c, &req)
	if err != nil {
		if c.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, respx.ResponseFail("service timeout", c.Err()))
			return
		}

		c.JSON(http.StatusInternalServerError, respx.ResponseFail("failed creating asssistant", err))
		return
	}

	c.JSON(http.StatusCreated, respx.ResponseSucceed("Assistant successfully created", dto.NewCreateAssistantResponse(resp)))
}

func (h *AuthHandler) SetPassword(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	// defer cancel()

	var req pbauth.SetAssistantPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid body", err))
		return
	}

	ctx := setLoginDataToContext(c)

	if _, err := h.assistantClient.SetAssistantPassword(ctx, &req); err != nil {
		if c.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, respx.ResponseFail("service timeout", c.Err()))
			return
		}

		c.JSON(http.StatusInternalServerError, respx.ResponseFail("failed updating assistant", err))
		return
	}

	c.JSON(http.StatusOK, respx.ResponseSucceed("Assistant successfully updated", nil))
}
