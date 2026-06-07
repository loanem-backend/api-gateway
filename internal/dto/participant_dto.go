package dto

import (
	"time"

	pbparticipant "github.com/loanem-backend/protos/pb/proto/services/participant/v1"
)

type ClassResponse struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Course    *withCourse `json:"course,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func PBClassToClassResponse(c *pbparticipant.Class) *ClassResponse {
	return &ClassResponse{
		ID:   int(c.GetId()),
		Name: c.GetName(),
		Course: &withCourse{
			ID:   int(c.GetCourseId()),
			Name: c.GetCourseName(),
			Year: int(c.GetCourseYear()),
		},
		CreatedAt: c.GetCreatedAt().AsTime(),
		UpdatedAt: c.GetUpdatedAt().AsTime(),
	}
}

type ClassesByCourseResponse struct {
	Course  *withCourse     `json:"course"`
	Classes []ClassResponse `json:"classes"`
}

func PBClassToClassByCourseResponse(c *pbparticipant.Class) *ClassResponse {
	return &ClassResponse{
		ID:        int(c.GetId()),
		Name:      c.GetName(),
		CreatedAt: c.GetCreatedAt().AsTime(),
		UpdatedAt: c.GetUpdatedAt().AsTime(),
	}
}

func GetClassesByCourseIDResponseToClassResponses(resp *pbparticipant.GetClassesByCourseIDResponse) *ClassesByCourseResponse {
	classes := resp.GetClasses()
	classCount := len(classes)

	if classCount < 1 {
		return &ClassesByCourseResponse{}
	}

	classResponses := make([]ClassResponse, classCount)

	for i, c := range classes {
		classResponses[i] = *PBClassToClassByCourseResponse(c)
	}

	return &ClassesByCourseResponse{
		Course: &withCourse{
			ID:   int(resp.GetClasses()[0].GetCourseId()),
			Name: resp.GetClasses()[0].GetCourseName(),
			Year: int(resp.GetClasses()[0].GetCourseYear()),
		},
		Classes: classResponses,
	}
}
