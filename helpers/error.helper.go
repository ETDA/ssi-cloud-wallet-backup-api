package helpers

import (
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"net/http"
)

func HTTPErrorToIError(res *core.RequestResponse) core.IError {
	ierr := &core.Error{}
	err := utils.MapToStruct(res.Data, ierr)
	if err != nil {
		return core.Error{
			Status:  http.StatusInternalServerError,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		}
	}
	ierr.Status = res.StatusCode
	return *ierr
}
