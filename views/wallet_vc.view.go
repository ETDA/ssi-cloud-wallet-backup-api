package views

import (
	"gitlab.finema.co/finema/etda/vc-wallet-api/models"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"time"
)

type WalletVC struct {
	ID           string     `json:"id"`
	CID          string     `json:"cid"`
	SchemaType   string     `json:"schema_type"`
	IssuanceDate *time.Time `json:"issuance_date"`
	JWT          string     `json:"jwt"`
	Issuer       string     `json:"issuer"`
	Holder       string     `json:"holder"`
	Status       string     `json:"status"`
}

func NewWalletVC(vc *models.VC, vcStatus *VCStatus) *WalletVC {
	view := &WalletVC{}
	_ = utils.Copy(view, vc)

	if vcStatus != nil {
		view.Status = utils.GetString(vcStatus.Status)
	}

	return view
}

func NewWalletVCs(vcs []models.VC, vcStatuses []VCStatus) []WalletVC {
	views := make([]WalletVC, len(vcs))

	for i, vc := range vcs {
		views[i] = *NewWalletVC(&vc, &VCStatus{})
	}

	return replaceStatusToVC(views, vcStatuses)
}

func replaceStatusToVC(vcWallets []WalletVC, vcStatuses []VCStatus) []WalletVC {
	views := make([]WalletVC, 0)

	for _, vpStatus := range vcStatuses {
		for _, vcWallet := range vcWallets {
			if vcWallet.CID == vpStatus.CID {
				vc := vcWallet
				vc.Status = utils.GetString(vpStatus.Status)
				views = append(views, vc)
				break
			}
		}
	}

	return views
}
