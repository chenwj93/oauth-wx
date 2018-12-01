package services

import (
	"oauth-wx/models"

	"encoding/json"
	"utils"
)

type RequestLogService struct {
	reqLog *models.RequestLog
}

func (r *RequestLogService) SaveLog(token interface{}, req string, param interface{}, resp []byte) {
	jParam, _ := json.Marshal(param)
	r.reqLog = &models.RequestLog{Token: utils.ParseString(token),
		Request: req,
		Param: string(jParam),
		Response: string(resp),
	}
	r.reqLog.Add()
}
