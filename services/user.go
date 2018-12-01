package services

import (
	"oauth-wx/models"
	"time"
	"utils"
	"cloud-config-client-go/conf"
)



type User struct {
	user *models.UserExt
}

func (u *User) GetModelUser() interface{} {
	return u.user.UserInfo
}

func (u *User) Check(openId interface{}) bool{
	var err error
	open_id := utils.ParseString(openId)
	if open_id == utils.EMPTY_STRING{
		return false
	}

	u.user, err = models.GetUser(open_id, conf.GetInteger("expiry.hour"))
	if err != nil || u.user == nil {
		return false
	}
	return true
}

func (u *User) GetToken(openId interface{}) (token string){
	token = u.user.GenerateToken()
	return
}

func (u *User) Decode(accessToken string) (t time.Time, err error) {
	u.user, err = models.DecodeToken(accessToken)
	if err == nil{
		t = time.Unix(u.user.Expiry, 0)
	}
	return
}
