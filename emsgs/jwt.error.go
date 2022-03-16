package emsgs

import (
	core "ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

var (
	JWTInValid = core.Error{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_JWT",
		Message: "JWT is not valid",
	}
)
