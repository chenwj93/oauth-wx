package models

import (
	"logs"
	"github.com/astaxie/beego/orm"
)

type ApiConfig struct {
	Api        string     `json:"api" orm:"column(api);pk"`
	Status     int8    `json:"status" orm:"column(status);"`
}

func (m *ApiConfig) TableName() string {
	return "api_config"
}

func init() {
	orm.RegisterModel(new(ApiConfig))
}

func GetApiList() (api map[string]bool, err error) {
	o := orm.NewOrm()
	var apis []ApiConfig
	num, err := o.Raw("select * from api_config where status = 1 and api is not null").QueryRows(&apis)
	if err != nil {
		return
	}
	api = make(map[string]bool, num)
	for i := range apis{
		if apis[i].Api != ""{
			logs.Info(apis[i].Api)
			api[apis[i].Api] = true
		}
	}
	return
}
