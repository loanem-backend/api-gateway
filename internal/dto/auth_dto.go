package dto

import (
	"time"

	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
)

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginRequestDTOToPB(req *LoginRequest) *pbauth.LoginRequest {
	return &pbauth.LoginRequest{
		Phone:    req.Phone,
		Password: req.Password,
	}
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func NewLoginResponse(token string) *LoginResponse {
	return &LoginResponse{
		AccessToken: token,
	}
}

type CreateAssistantResponse struct {
	ID int `json:"id"`
}

func NewCreateAssistantResponse(resp *pbauth.CreateAssistantResponse) *CreateAssistantResponse {
	return &CreateAssistantResponse{
		ID: int(resp.GetId()),
	}
}

type AssistantResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Period    int       `json:"period"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func PBAssistantToAssistantResponse(a *pbauth.Assistant) *AssistantResponse {
	return &AssistantResponse{
		ID:        int(a.GetId()),
		Name:      a.GetName(),
		Phone:     a.GetPhone(),
		Period:    int(a.GetPeriod()),
		IsActive:  a.GetActive(),
		CreatedAt: a.GetCreatedAt().AsTime(),
		UpdatedAt: a.GetUpdatedAt().AsTime(),
	}
}

func GetAssistantsResponse(resp *pbauth.GetActiveAssistantsResponse) []AssistantResponse {
	assistants := resp.GetAssistants()
	responses := make([]AssistantResponse, len(assistants))

	for i, a := range assistants {
		responses[i] = *PBAssistantToAssistantResponse(a)
	}

	return responses
}
