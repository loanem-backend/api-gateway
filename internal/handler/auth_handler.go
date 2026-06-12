package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

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
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid body", err))
		return
	}

	resp, err := h.authClient.Login(c.Request.Context(), dto.LoginRequestDTOToPB(&req))
	if err != nil {
		handleErrorFromClient(c, err)
		return
	}

	c.SetCookie(
		cookieNameRefreshToken, resp.GetRefreshToken(),
		int(resp.GetRefreshExpirationHour())*3600, "/", "", false, true,
	)

	c.JSON(http.StatusCreated, respx.ResponseSucceed(messageLoginSucceed, dto.NewLoginResponse(resp.GetAccessToken())))
}

const cookieNameRefreshToken = "refresh_token"

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

func getAuthorization(c *gin.Context) (string, error) {
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

func (h *AuthHandler) Logout(c *gin.Context) {
	accessToken, err := getAuthorization(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, respx.ResponseFail("unauthorized", err))
	}

	refreshToken, err := c.Cookie(cookieNameRefreshToken)
	if err != nil {
		refreshToken = ""
	}

	ctx := setLoginDataToContext(c)

	if _, err := h.authClient.Logout(ctx, &pbauth.LogoutRequest{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}); err != nil {
		if c.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, respx.ResponseFail("service timeout", c.Err()))
			return
		}

		c.JSON(http.StatusInternalServerError, respx.ResponseFail("failed logging out", err))
		return
	}

	c.JSON(http.StatusOK, respx.ResponseSucceed("Successfully logged out", nil))
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie(cookieNameRefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, respx.ResponseFail("refresh token not found", errors.New("missing cookie")))
		return
	}

	resp, err := h.authClient.RefreshToken(c, &pbauth.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		if c.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, respx.ResponseFail("service timeout", c.Err()))
			return
		}

		c.JSON(http.StatusInternalServerError, respx.ResponseFail("failed generating token", err))
		return
	}

	c.SetCookie(
		cookieNameRefreshToken, resp.GetRefreshToken(),
		int(resp.GetRefreshExpirationHour())*3600, "/", "", false, true,
	)

	c.JSON(http.StatusCreated, respx.ResponseSucceed("Successfully logged in", dto.NewLoginResponse(resp.GetAccessToken())))
}
