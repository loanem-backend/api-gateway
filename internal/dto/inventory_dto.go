package dto

import (
	"time"

	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
)

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

type InstrumentResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func PBInstrumentToInstrumentResponse(i *pbinventory.Instrument) *InstrumentResponse {
	return &InstrumentResponse{
		ID:        int(i.GetId()),
		Name:      i.GetName(),
		Picture:   i.GetPicture(),
		CreatedAt: i.GetCreatedAt().AsTime(),
		UpdatedAt: i.GetUpdatedAt().AsTime(),
	}
}

func GetAllInstrumentsResponseToInstrumentResponses(resp *pbinventory.GetAllInstrumentsResponse) []InstrumentResponse {
	instruments := resp.GetInstruments()

	responses := make([]InstrumentResponse, len(instruments))

	for idx, i := range instruments {
		responses[idx] = *PBInstrumentToInstrumentResponse(i)
	}

	return responses
}
