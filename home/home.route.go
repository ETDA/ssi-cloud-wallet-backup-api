package home

import (
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/vc-wallet-api/middlewares"
	core "ssi-gitlab.teda.th/ssi/core"
)

func NewHomeHTTPHandler(r *echo.Echo) {
	home := &HomeController{}
	r.GET("/", core.WithHTTPContext(home.Get))
	r.GET("wallet/:did", core.WithHTTPContext(home.Find), middlewares.VerifySignatureMiddleware)
	r.POST("wallet", core.WithHTTPContext(home.Create), middlewares.VerifySignatureMiddleware)
	r.GET("/wallet/:did/vcs/:cid", core.WithHTTPContext(home.CheckVC), middlewares.VerifySignatureMiddleware)
	r.POST("/wallet/:did/vcs", core.WithHTTPContext(home.AddVC), middlewares.VerifySignatureMiddleware)
	r.GET("/wallet/:did/vcs", core.WithHTTPContext(home.VCPagination), middlewares.VerifySignatureMiddleware)
}
