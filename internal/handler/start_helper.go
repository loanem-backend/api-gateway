package handler

import (
	"strings"

	"github.com/loanem-backend/api-gateway/internal/config"
)

func serviceAddr(idx serviceIdx) string {
	sName := serviceName(idx)
	sNameUpper := strings.ToUpper(sName)

	addr := config.GetEnv(sNameUpper+"_SERV_ADDR", sName+"-service")

	port := config.GetEnv(sNameUpper+"_SERV_PORT", "50054")

	return addr + ":" + port
}

type serviceIdx int

const (
	auth serviceIdx = iota
	course
	inventory
	participant
)

func serviceName(idx serviceIdx) string {
	switch idx {
	case auth:
		return "auth"
	case course:
		return "course"
	case inventory:
		return "inventory"
	case participant:
		return "participant"
	default:
		return ""
	}
}
