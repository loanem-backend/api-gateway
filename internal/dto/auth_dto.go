package dto

import (
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
