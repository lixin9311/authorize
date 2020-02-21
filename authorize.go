package main

import (
	"github.com/lixin9311/authorize/model"
	"github.com/lixin9311/osin"

	"github.com/labstack/echo"
)

type AutorizeHandler struct {
	oauthServer *osin.Server
}

func (h *AutorizeHandler) Authorize(c echo.Context) (err error) {
	resp := h.oauthServer.NewResponse()
	defer resp.Close()

	if ar := h.oauthServer.HandleAuthorizeRequest(resp, c); ar != nil {
		u := model.User{}
		c.Bind(&u)
		if u.Username == "test" && u.Password == "test" {
			ar.Authorized = true
			ar.UserData = u
		}
		h.oauthServer.FinishAuthorizeRequest(resp, c, ar)
	}
	return osin.OutputJSON(resp, c)
}

func (h *AutorizeHandler) Token(c echo.Context) (err error) {
	resp := h.oauthServer.NewResponse()
	defer resp.Close()

	if ar := h.oauthServer.HandleAccessRequest(resp, c); ar != nil {
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		case osin.PASSWORD:
			if ar.Username == "test" && ar.Password == "test" {
				ar.Authorized = true
			}
		case osin.CLIENT_CREDENTIALS:
			ar.Authorized = true
		}
		h.oauthServer.FinishAccessRequest(resp, c, ar)
	}
	return osin.OutputJSON(resp, c)
}
