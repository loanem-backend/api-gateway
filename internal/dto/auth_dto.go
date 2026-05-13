package dto

import pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"

type LoginResponse struct {
	Token string `json:"token"`
}

func NewLoginResponse(resp *pbauth.LoginResponse) *LoginResponse {
	return &LoginResponse{
		Token: resp.GetToken(),
	}
}
