package main

import (
	"go-fx/resources"
	"os"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterModel(new(resources.Job))
	orm.RegisterModel(new(resources.Run))
	orm.RegisterDataBase("default", "mysql", os.Getenv("MYSQL_URL"))
}

func main() {
	orm.RunSyncdb("default", false, false)
	web.Get("/run/:name", resources.RunJob)
	web.Run("0.0.0.0:8080")
}
