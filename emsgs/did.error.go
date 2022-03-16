package emsgs

import (
	core "ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

var DuplicatedDIDAddress = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "DUPLICATED_DID_ADDRESS",
	Message: "did address already used",
}

var InvalidDIDAddress = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "INVALID_DID_ADDRESS",
	Message: "did address is invalid format or not registered",
}
