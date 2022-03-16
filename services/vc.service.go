package services

import (
	"fmt"
	"gitlab.finema.co/finema/etda/vc-wallet-api/consts"
	"gitlab.finema.co/finema/etda/vc-wallet-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"strings"
)

type IVCService interface {
	Find(cid string) (*views.VCStatus, core.IError)
	FindMultiple(cids []string, options *VCStatusFindMultipleOptions) ([]views.VCStatus, core.IError)
}

type vcService struct {
	ctx core.IContext
}

func NewVCService(ctx core.IContext) IVCService {
	return &vcService{ctx: ctx}
}

type VCStatusFindMultipleOptions struct {
	Status string
}

func (s *vcService) Find(cid string) (*views.VCStatus, core.IError) {
	vc := &views.VCStatus{}

	ierr := core.RequesterToStruct(vc, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(fmt.Sprintf("/vc/status/%s", cid), &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVVCStatusServiceBaseURL),
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return vc, nil
}

func (s *vcService) FindMultiple(cids []string, options *VCStatusFindMultipleOptions) ([]views.VCStatus, core.IError) {
	vcs := make([]views.VCStatus, 0)

	ierr := core.RequesterToStruct(&vcs, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(fmt.Sprintf("/vc/status?cid=%s", strings.Join(cids, ",")), &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVVCStatusServiceBaseURL),
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if options != nil {
		if options.Status != "" {
			tmp := vcs
			vcs = make([]views.VCStatus, 0)

			for _, vc := range tmp {
				if utils.GetString(vc.Status) == options.Status {
					vcs = append(vcs, vc)
				}
			}
		}
	}

	return vcs, nil
}
