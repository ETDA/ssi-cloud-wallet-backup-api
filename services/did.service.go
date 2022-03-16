package services

import (
	"fmt"
	"gitlab.finema.co/finema/etda/vc-wallet-api/consts"
	"gitlab.finema.co/finema/etda/vc-wallet-api/helpers"
	"gitlab.finema.co/finema/etda/vc-wallet-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type IDIDService interface {
	Find(address string) (*models.DIDDocument, core.IError)
}
type didService struct {
	ctx core.IContext
}

func NewDIDService(ctx core.IContext) IDIDService {
	return &didService{ctx: ctx}
}

func (s *didService) Find(address string) (*models.DIDDocument, core.IError) {
	res, err := s.ctx.Requester().Get(fmt.Sprintf("/did/%s/document/latest", address),
		&core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVDIDServiceBaseURL),
		})
	if res == nil {
		return nil, s.ctx.NewError(errmsgs.InternalServerError, errmsgs.InternalServerError)
	}

	if err != nil {
		ierr := helpers.HTTPErrorToIError(res)
		return nil, s.ctx.NewError(ierr, ierr)
	}

	result := &models.DIDDocument{}
	err = utils.MapToStruct(res.Data, result)
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.InternalServerError)
	}

	return result, nil
}
