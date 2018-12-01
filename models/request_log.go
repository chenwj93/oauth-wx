package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type RequestLog struct {
	Id         int       `json:"id" orm:"column(id);pk"`
	Token      string    `json:"token" orm:"column(token);size(40);null"`
	Request    string    `json:"request" orm:"column(request);size(50);null"`
	Param      string    `json:"param" orm:"column(param);size(255);null"`
	Response   string    `json:"response" orm:"column(response);size(255);null"`
	CreateTime time.Time `json:"-" orm:"column(create_time);type(timestamp);null"`
}

func (m *RequestLog) TableName() string {
	return "request_log"
}

func init() {
	orm.RegisterModel(new(RequestLog))
}

func (m *RequestLog) Add() (err error) {
	m.CreateTime = time.Now()
	o := orm.NewOrm()
	_, err = o.Insert(m)
	return
}

func GetRequestLogById(id int) (v *RequestLog, err error) {
	o := orm.NewOrm()
	v = &RequestLog{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func UpdateRequestLogById(m *RequestLog) (err error) {
	o := orm.NewOrm()
	v := RequestLog{Id: m.Id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func DeleteRequestLog(id int) (err error) {
	o := orm.NewOrm()
	v := RequestLog{Id: id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&RequestLog{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}


