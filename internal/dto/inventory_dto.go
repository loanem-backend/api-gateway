package dto

import pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"

type CreateInstrumentResponse struct {
	ID int `json:"id"`
}

func NewCreateInstrumentResponse(resp *pbinventory.AddInstrumentResponse) *CreateInstrumentResponse {
	return &CreateInstrumentResponse{
		ID: int(resp.GetId()),
	}
}
