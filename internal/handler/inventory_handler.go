package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loanem-backend/api-gateway/internal/dto"
	"github.com/loanem-backend/api-gateway/pkg/respx"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
)

type InventoryHandler struct {
	instrumentClient pbinventory.InstrumentServiceClient
	toolkitClient    pbinventory.ToolkitServiceClient
}

func NewInventoryHandler(ic pbinventory.InstrumentServiceClient, tc pbinventory.ToolkitServiceClient) *InventoryHandler {
	return &InventoryHandler{
		instrumentClient: ic,
		toolkitClient:    tc,
	}
}

func (h *InventoryHandler) CreateInstrument(c *gin.Context) {
	var req pbinventory.AddInstrumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid body", err))
		return
	}

	ctx := setLoginDataToContext(c)

	resp, err := h.instrumentClient.AddInstrument(ctx, &req)
	if err != nil {
		if c.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, respx.ResponseFail("service timeout", c.Err()))
			return
		}

		c.JSON(http.StatusInternalServerError, respx.ResponseFail("failed creating instrument", err))
		return
	}

	c.JSON(http.StatusCreated, respx.ResponseSucceed("Instrument successfully created", dto.NewCreateInstrumentResponse(resp)))
}

func (h *InventoryHandler) CreateToolkit(c *gin.Context) {
	var req pbinventory.AddToolkitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid body", err))
		return
	}

	ctx := setLoginDataToContext(c)

	resp, err := h.toolkitClient.AddToolkit(ctx, &req)
	if err != nil {
		if c.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, respx.ResponseFail("service timeout", c.Err()))
			return
		}

		c.JSON(http.StatusInternalServerError, respx.ResponseFail("failed creating toolkit", err))
		return
	}

	c.JSON(http.StatusCreated, respx.ResponseSucceed("Toolkit successfully created", dto.NewCreateToolkitResponse(resp)))
}
