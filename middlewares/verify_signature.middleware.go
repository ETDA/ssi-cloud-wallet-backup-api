package middlewares

import (
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/vc-wallet-api/consts"
	"gitlab.finema.co/finema/etda/vc-wallet-api/emsgs"
	"gitlab.finema.co/finema/etda/vc-wallet-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"net/http"
)

type payload struct {
	core.BaseValidator
	Message *string `json:"message"`
}

func (r payload) Valid(ctx core.IContext) core.IError {
	if r.Must(r.IsStrRequired(r.Message, "message")) {
		r.Must(r.IsBase64(r.Message, "message"))
	}

	return r.Error()
}

type OperationPayload struct {
	core.BaseValidator
	Operation  *string `json:"operation"`
	DIDAddress *string `json:"did_address"`
}

func (r OperationPayload) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Operation, "operation"))
	r.Must(r.IsStrRequired(r.DIDAddress, "did_address"))

	return r.Error()
}

func VerifySignatureMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(core.IHTTPContext)
		if cc.GetSignature() == "" {
			return c.JSON(http.StatusBadRequest, core.NewValidatorFields(core.RequiredM("x-signature")))
		}

		var isSigValid = false
		var didAddress string
		var message string
		if cc.Request().Method == http.MethodPost {
			payloadData := &payload{}
			if err := cc.BindWithValidate(payloadData); err != nil {
				return c.JSON(err.GetStatus(), err.JSON())
			}

			message = utils.GetString(payloadData.Message)
			jsonString, err := utils.Base64Decode(message) // decode failed
			if err != nil {
				return c.JSON(errmsgs.BadRequest.GetStatus(), errmsgs.BadRequest.JSON())
			}

			messagePayload := &OperationPayload{}
			err = utils.JSONParse([]byte(jsonString), messagePayload) // unmarshall failed
			if err != nil {
				return c.JSON(errmsgs.BadRequest.GetStatus(), errmsgs.BadRequest.JSON())
			}
			if ierr := messagePayload.Valid(cc); ierr != nil {
				return c.JSON(ierr.GetStatus(), ierr.JSON())
			}

			didAddress = utils.GetString(messagePayload.DIDAddress)
			c.Set(consts.ContextKeyDIDAddress, didAddress) // only set did on operational request only
		}

		if cc.Request().Method == http.MethodGet {
			didAddress = cc.Param("did")
			message = didAddress
		}

		didService := services.NewDIDService(cc)
		didDocument, ierr := didService.Find(didAddress)

		if errmsgs.IsNotFoundError(ierr) {
			return c.JSON(emsgs.InvalidDIDAddress.GetStatus(), emsgs.InvalidDIDAddress.JSON())
		}

		if ierr != nil {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}

		for _, verificationMethod := range didDocument.VerificationMethod {
			valid, _ := utils.VerifySignature(verificationMethod.PublicKeyPem, cc.GetSignature(), message)
			if valid {
				isSigValid = valid
				break
			}
		}

		if !isSigValid {
			return c.JSON(errmsgs.SignatureInValid.GetStatus(), errmsgs.SignatureInValid.JSON())
		}

		c.Set(consts.ContextKeyMessage, message)
		return next(c)
	}
}
