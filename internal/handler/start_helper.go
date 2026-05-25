package handler

import (
	"strings"

	"github.com/loanem-backend/api-gateway/internal/config"
)

func serviceAddr(idx serviceIdx) string {
	sName := serviceName(idx)
	sNameUpper := strings.ToUpper(sName)

	addr := config.GetEnv(sNameUpper+"_SERV_ADDR", sName+"-service")

	port := config.GetEnv(sNameUpper+"_SERV_PORT", "50051")

	return addr + ":" + port
}

type serviceIdx int

const (
	auth serviceIdx = iota
	inventory
)

func serviceName(idx serviceIdx) string {
	switch idx {
	case auth:
		return "auth"
	case inventory:
		return "inventory"
	default:
		return ""
	}
}
