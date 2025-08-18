package resources

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	_ "os/exec"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web/context"
)

type Job struct {
	ID    int    `orm:"auto"`
	Name  string `orm:"column(name)"`
	Image string `orm:"column(image)"`
	Code  string `orm:"column(code)"`
	Runs  []*Run `orm:"reverse(many)"`
}

func (j *Job) CommandToExecute(ctx *context.Context) []string {
	switch j.Image {
	case "php":
		return []string{"php", "-r", j.Hydrate(ctx)}
	case "node":
		return []string{"node", "-e", j.Code}
	}

	return []string{}
}

func (j *Job) Hydrate(ctx *context.Context) string {
	header, err := json.Marshal(ctx.Request.Header)
	query, err := json.Marshal(ctx.Input.Params())
	if err != nil {
		fmt.Println(err)
	}
	return "class Request{public $headers;public $query;}" + j.Code + "$request = new Request();$request->headers = json_decode('" + string(header) + "');$request->query = json_encode('" + string(query) + "'); echo main($request);"
}

type Run struct {
	ID          int    `orm:"auto"`
	Job         *Job   `orm:"rel(fk)"`
	DateTime    string `orm:"column(datetime)"`
	TimeElapsed int64  `orm:"column(time_elapsed)"`
}

func RunJob(ctx *context.Context) {
	name := ctx.Input.Param(":name")
	fmt.Println("attempting to run " + name)

	var job Job
	o := orm.NewOrm()
	o.QueryTable((*Job)(nil)).Filter("name", name).One(&job)

	startTime := time.Now()
	cmdArgs := job.CommandToExecute(ctx)
	cmd := exec.Command("docker", append([]string{"run", "--rm", job.Image}, cmdArgs...)...)
	fmt.Println(cmd.Args)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Befehlsausführung fehlgeschlagen: %v", err)
	}

	ctx.WriteString(string(output))
	run := new(Run)
	run.Job = &job
	run.DateTime = time.Now().Format("2006-01-02 15:04:05")
	run.TimeElapsed = time.Since(startTime).Milliseconds()
	o.Insert(run)
	fmt.Println("done in " + strconv.FormatInt(time.Since(startTime).Milliseconds(), 10) + " ms")
}
