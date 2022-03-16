package requests

import (
	"gitlab.finema.co/finema/etda/vc-wallet-api/consts"
	core "ssi-gitlab.teda.th/ssi/core"
)

type WalletVCCreate struct {
	core.BaseValidator
	Operation  *string `json:"operation"`
	DIDAddress *string `json:"did_address"`
	JWT        *string `json:"jwt"`
}

func (r WalletVCCreate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrIn(r.Operation, consts.OperationWalletVCAdd, "operation"))
	r.Must(r.IsStrRequired(r.DIDAddress, "did_address"))
	r.Must(r.IsStrRequired(r.JWT, "jwt"))

	return r.Error()
}
