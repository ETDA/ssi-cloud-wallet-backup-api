package home

import (
	"net/http"

	"gitlab.finema.co/finema/etda/vc-wallet-api/consts"
	"gitlab.finema.co/finema/etda/vc-wallet-api/requests"
	"gitlab.finema.co/finema/etda/vc-wallet-api/services"
	"gitlab.finema.co/finema/etda/vc-wallet-api/views"
	"ssi-gitlab.teda.th/ssi/core/utils"

	core "ssi-gitlab.teda.th/ssi/core"
)

type HomeController struct{}

func (n *HomeController) Get(c core.IHTTPContext) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "I am VC wallet API",
	})
}

func (n *HomeController) Find(c core.IHTTPContext) error {
	service := services.NewWalletService(c)
	ierr := service.CheckWallet(c.Param("did"))
	if ierr != nil {
		return c.JSON(http.StatusOK, core.Map{
			"is_exists": false,
		})
	}
	return c.JSON(http.StatusOK, core.Map{
		"is_exists": true,
	})
}
func (n *HomeController) Create(c core.IHTTPContext) error {
	input := &requests.WalletCreate{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	service := services.NewWalletService(c)
	wallet, ierr := service.CreateWallet(&services.WalletCreatePayload{
		DIDAddress: utils.GetString(input.DIDAddress),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, wallet)
}

func (n *HomeController) CheckVC(c core.IHTTPContext) error {
	service := services.NewWalletService(c)
	ierr := service.CheckVC(c.Param("cid"))
	if ierr != nil {
		return c.JSON(http.StatusOK, core.Map{
			"is_exists": false,
		})
	}
	return c.JSON(http.StatusOK, core.Map{
		"is_exists": true,
	})
}
func (n *HomeController) AddVC(c core.IHTTPContext) error {
	input := &requests.WalletVCCreate{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	service := services.NewWalletService(c)
	ierr := service.AddVC(c.Param("did"), &services.WalletAddVCPayload{
		JWT:      utils.GetString(input.JWT),
		Operator: c.Get(consts.ContextKeyDIDAddress).(string),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.Map{"result": "success"})
}

func (n *HomeController) VCPagination(c core.IHTTPContext) error {
	service := services.NewWalletService(c)
	items, pageRes, ierr := service.VCPagination(c.Param("did"), c.GetPageOptions(), &services.VCPaginationOptions{
		SchemaType:   c.QueryParam("type"),
		IssuanceDate: c.QueryParam("issuance_date"),
		Holder:       c.QueryParam("holder"),
		Issuer:       c.QueryParam("issuer"),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	cids := make([]string, 0)
	for _, item := range items {
		cids = append(cids, item.CID)
	}

	vcService := services.NewVCService(c)
	vcStatuses, ierr := vcService.FindMultiple(cids, &services.VCStatusFindMultipleOptions{
		Status: c.QueryParam("status"),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.NewPagination(views.NewWalletVCs(items, vcStatuses), pageRes))
}
