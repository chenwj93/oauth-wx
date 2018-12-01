package controllers

import (
	"utils"
	"oauth-wx/services"
	"context"
	"logs"
	"time"
)

type ILogin interface {
	Login() (*utils.Response, error)
}

type IAuthenticate interface {
	Authenticate() (*utils.Response, error)
}

type IAuthModify interface {
	AuthHadModified() (*utils.Response, error)
}

type Access struct {
	Ctx    context.Context
	reqLog services.RequestLogService
	user   services.User
	global services.GlobalService
}

func (a *Access) Login() (*utils.Response, error) {
	paramInput := a.GetParamMap()
	ok := a.user.Check(paramInput["openId"])
	if !ok {
		return utils.ConcatFailed(nil, "请先注册会员！")
	}
	act := a.user.GetToken(paramInput["openId"])

	return utils.ConcatSuccess(act)
}

func (a *Access) Authenticate() (*utils.Response, error) {
	act := utils.ParseString(a.Ctx.Value("access_token"))
	if act == utils.EMPTY_STRING {
		return utils.ConcatDeny()
	}
	expiry, err := a.user.Decode(act)
	if err != nil || expiry.Before(time.Now()) {
		logs.Error(err)
		logs.Info(expiry)
		return utils.ConcatDeny()
	}
	paramInput := a.GetParamMap()
	ok := a.global.Auth(paramInput["operate"])
	if !ok {
		return utils.ConcatDeny()
	}
	return utils.ConcatSuccess(a.user.GetModelUser())
}

func (a *Access) AuthHadModified() (*utils.Response, error) {
	a.global.MarkAuthModify()
	return utils.ConcatSuccess()
}

func (a *Access) GetParamMap() (m map[string]interface{}) {
	m, ok := a.Ctx.Value("param").(map[string]interface{})
	if !ok {
		logs.Error("param asset error :", a.Ctx.Value("param"))
	}
	return
}
