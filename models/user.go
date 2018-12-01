package models

import (
	"fmt"
	"time"
	"encoding/json"
	"logs"

	"github.com/astaxie/beego/orm"
	"utils"
	"cloud-config-client-go/conf"

	"oauth-wx/utils-local"
	"oauth-wx/constant"
)

type User struct {
	OpenId     string `json:"openId" orm:"column(open_id);pk"`
	UserName   string `json:"userName" orm:"column(user_name);size(64);null"`
	UserId     int64  `json:"userId" orm:"column(user_id);null"`
	Phone      string `json:"phone" orm:"column(phone);null"`
	ClientId   string `json:"clientId" orm:"column(client_id);"`
	UpdateTime string `json:"-" orm:"column(update_time);size(19);null"`
	Deleted    int8   `json:"-" orm:"column(deleted);null"`
}

type UserExt struct {
	UserInfo User
	Expiry   int64 `json:"expiry"`
}

func (m *User) TableName() string {
	return "user"
}

func init() {
	orm.RegisterModel(new(User))
}

func (m *User) init(param interface{}) {
	utils.ParseStruct(param, m)
	m.UpdateTime = utils.Now()
	m.Deleted = utils.ZERO_I
}

func (m *User) Add() (err error) {
	o := orm.NewOrm()
	_, err = o.Insert(m)
	return
}

func GetUser(openId string, expiry int) (*UserExt, error) {
	o := orm.NewOrm()
	var v UserExt
	err := o.Raw("select * from user where deleted = 0 and open_id = ?", openId).QueryRow(&v.UserInfo)
	ifExpiry := false
	now := time.Now().Add(-time.Hour * time.Duration(expiry)).Format(utils.TIME_FORMAT_1)
	if err == nil && v.UserInfo.OpenId != utils.EMPTY_STRING && v.UserInfo.UpdateTime < now {
		ifExpiry = true
	}
	if err != nil || v.UserInfo.OpenId == utils.EMPTY_STRING || ifExpiry {
		uMap, err := utils.ServiceCall(conf.GetString("member.center.app"), "getUserInfo", "get", map[string]interface{}{"openId": openId}, map[string]interface{}{"clientId": "-1"})
		if err == nil {
			logs.Error(err)
			v.UserInfo.init(uMap)
			if ifExpiry {
				err = v.UserInfo.Update([]string{"UpdateTime"})
			} else {
				err = v.UserInfo.Add()
			}
		}
	}

	return &v, err
}

func (m *User) Update(fields []string) (err error) {
	o := orm.NewOrm()
	if err = o.Read(m); err == nil {
		var num int64
		if num, err = o.Update(m, fields...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func (m *UserExt) GenerateToken() string {
	m.Expiry = time.Now().Unix() + 60*60*int64(conf.GetInteger("expiry.hour"))
	j, _ := json.Marshal(*m)
	result := utils_local.AesEncrypt(string(j), constant.Private_key, constant.Slash)
	return result
}

func DecodeToken(token string) (v *UserExt, err error) {
	origData, err := utils_local.AesDecrypt(token, constant.Private_key, constant.Slash)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	var u UserExt
	err = json.Unmarshal([]byte(origData), &u)
	v = &u
	return
}
