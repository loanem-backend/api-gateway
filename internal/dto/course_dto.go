package dto

import pbcourse "github.com/loanem-backend/protos/pb/proto/services/course/v1"

type CreateCourseResponse struct {
	ID int `json:"id"`
}

func NewCreateCourseResponse(req *pbcourse.AddCourseResponse) *CreateCourseResponse {
	return &CreateCourseResponse{
		ID: int(req.GetId()),
	}
}
