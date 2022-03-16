package emsgs

import (
	core "ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

var WalletNotFound = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "WALLET_NOT_FOUND",
	Message: "wallet is not found",
}

var DuplicatedWallet = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "DUPLICATED_WALLET",
	Message: "wallet already exists",
}
