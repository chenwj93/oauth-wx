package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"oauth-wx/routers"
	"net/http"
	"cloud-config-client-go/conf"
	"flag"
	"logs"
	"utils"
	"log"
	"oauth-wx/constant"
	"oauth-wx/models"
)

func init() {
	var PROFILE string
	flag.StringVar(&PROFILE,"profile","test","execute environment")
	flag.Parse()
	fmt.Println(PROFILE)
	conf.Load(constant.EUREKA_SERVER, constant.CONFIG_SERVER, constant.APP_NAME, PROFILE, constant.BRANCH, false)
	// conf.LoadLocalProperties("conf.properties")

	//database connection
	mysql_user := conf.GetString("mysql.username")
	mysql_pass := conf.GetString("mysql.password")
	mysql_urls := conf.GetString("mysql.url")
	mysql_port := conf.GetString("mysql.port")
	mysql_name := conf.GetString("mysql.conn.name")

	// timezone
	orm.DefaultTimeLoc = time.Local
	orm.RegisterDataBase("default", "mysql", mysql_user+":"+mysql_pass+"@tcp("+mysql_urls+":"+mysql_port+")/"+mysql_name+ "?charset=utf8&loc=Asia%2FShanghai")
	orm.Debug = true
	runtime.GOMAXPROCS(runtime.NumCPU())

	//machines :=[]string{constant.EUREKA_SERVER}
	//utils.StartEureka(constant.APP_NAME, conf.GetInteger("app.port"), nil, machines)

	//logs.InitSingleFile("log/log.log", orm.DebugLog.Logger)
	logs.InitTimeFile("log/saas/oauth-wx/", time.Hour * 24, []*log.Logger{orm.DebugLog.Logger, utils.ULog.Logger}, logs.DEBUG)
	models.ApiCacheInstance.Init()
}

// create by  chenwj on 2018-07-25
func main() {
	http.HandleFunc("/", routers.GetHandler())
	fmt.Printf("server at :%s\n", conf.GetString("app.port"))
	http.ListenAndServe(":" + conf.GetString("app.port"), nil)
}