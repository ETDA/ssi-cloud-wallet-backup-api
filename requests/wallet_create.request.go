package requests

import (
	"gitlab.finema.co/finema/etda/vc-wallet-api/consts"
	core "ssi-gitlab.teda.th/ssi/core"
)

type WalletCreate struct {
	core.BaseValidator
	Operation  *string `json:"operation"`
	DIDAddress *string `json:"did_address"`
}

func (r WalletCreate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrIn(r.Operation, consts.OperationWalletCreate, "operation"))
	r.Must(r.IsStrRequired(r.DIDAddress, "did_address"))

	return r.Error()
}
