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
	inventoryClient pbinventory.InventoryServiceClient
}

func NewInventoryHandler(ic pbinventory.InventoryServiceClient) *InventoryHandler {
	return &InventoryHandler{
		inventoryClient: ic,
	}
}

func (h *InventoryHandler) Create(c *gin.Context) {
	var req pbinventory.AddInstrumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respx.ResponseFail("invalid body", err))
		return
	}

	ctx := setLoginDataToContext(c)

	resp, err := h.inventoryClient.AddInstrument(ctx, &req)
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
