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

type CreateToolkitResponse struct {
	ID int `json:"id"`
}

func NewCreateToolkitResponse(resp *pbinventory.AddToolkitResponse) *CreateToolkitResponse {
	return &CreateToolkitResponse{
		ID: int(resp.GetId()),
	}
}
