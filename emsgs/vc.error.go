package emsgs

import (
	core "ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

var IssuerMismatched = core.Error{
	Status:  http.StatusForbidden,
	Code:    "DID_MISMATCH",
	Message: "only the issuer and holder are allowed to add vc to wallet",
}

var VCHolderInvalid = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "VC_HOLDER_INVALID",
	Message: "holder is not match with your did address",
}

var VCIssuerInvalid = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "VC_ISSUER_INVALID",
	Message: "issuer is not match with your did address",
}
