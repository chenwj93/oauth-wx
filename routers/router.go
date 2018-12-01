package routers

import (
	"oauth-wx/controllers"
	"utils"
	"net/http"
	"oauth-wx/constant"
	"logs"
	"context"
	"runtime/debug"
	"oauth-wx/services"
)

type BusinessServiceImpl struct {
	conLog services.RequestLogService
	auth   controllers.IAuthenticate
	login  controllers.ILogin
	authMd controllers.IAuthModify
}

// create by chenwj on 2018-07-25 2018-03-07
func (msi *BusinessServiceImpl) Handle(operation string, paramInput map[string]interface{}) (r *utils.Response, err error) {
	defer func() {
		if re := recover(); re != nil {
			logs.Fatal(re)
			logs.Fatal(string(debug.Stack()))
			r, err = utils.ConcatFailed(nil)
		}
	}()
	logs.Info("-->param:", utils.SimplifyStringMap(paramInput, 300))

	token := paramInput["access_token"]
	delete(paramInput, "access_token")
	ctx := context.WithValue(context.Background(), "access_token", token)
	ctx = context.WithValue(ctx, "param", paramInput)
	switch operation {
	case "login":
		msi.login = &controllers.Access{Ctx: ctx}
		r, err = msi.login.Login()
	case "authenticate":
		msi.auth = &controllers.Access{Ctx: ctx}
		r, err = msi.auth.Authenticate()
	case "authHadModified":
		msi.authMd = &controllers.Access{Ctx: ctx}
		r, err = msi.authMd.AuthHadModified()
	default:
		r, err = utils.ConcatNotFound()
	}
	go msi.conLog.SaveLog(token, operation, paramInput, r.Json)
	return
}



func GetHandler() http.HandlerFunc {
	f := func() utils.RouterInterface {
		return &BusinessServiceImpl{}
	}
	handler := utils.Handler{constant.ROOT_PATH, nil, f}
	return handler.Handle
}
