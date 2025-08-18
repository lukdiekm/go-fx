package main

import (
	"fmt"
	"go-fx/resources"
	"os"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterModel(new(resources.Job))
	orm.RegisterModel(new(resources.Run))
	fmt.Println("teeest" + os.Getenv("MYSQL_URL"))
	orm.RegisterDataBase("default", "mysql", os.Getenv("MYSQL_URL"))
}

func main() {
	orm.RunSyncdb("default", false, false)
	web.Get("/run/:name", resources.RunJob)
	web.Run("127.0.0.1:8080")
}
