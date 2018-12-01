package services

import (
	"utils"
	"oauth-wx/models"
)

type GlobalService struct {
}

func (g *GlobalService) Auth(operate interface{}) bool {
	operateStr := utils.ParseString(operate)
	if operateStr == utils.EMPTY_STRING { // || !strings.HasPrefix(operateStr, conf.GetString("api.prefix")){
		return false
	}
	ifAllow := models.ApiCacheInstance.GetApi(operateStr)
	return ifAllow
}
func (g *GlobalService) MarkAuthModify() {
	models.ApiCacheInstance.SetUpdated(true)
}
