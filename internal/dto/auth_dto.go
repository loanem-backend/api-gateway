package dto

import pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func NewLoginResponse(resp *pbauth.LoginResponse) *LoginResponse {
	return &LoginResponse{
		AccessToken: resp.GetAccessToken(),
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
