package dto

import (
	"time"

	pbcourse "github.com/loanem-backend/protos/pb/proto/services/course/v1"
)

type CreateCourseResponse struct {
	ID int `json:"id"`
}

func NewCreateCourseResponse(req *pbcourse.AddCourseResponse) *CreateCourseResponse {
	return &CreateCourseResponse{
		ID: int(req.GetId()),
	}
}

type CourseResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type withCourse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Year int    `json:"year"`
}
