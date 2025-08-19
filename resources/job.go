package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-fx/templates"
	"html"
	"html/template"
	"log"
	"os"
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
	Token string `orm:"column(token)"`
}

type Run struct {
	ID          int    `orm:"auto"`
	Job         *Job   `orm:"rel(fk)"`
	DateTime    string `orm:"column(datetime)"`
	TimeElapsed int64  `orm:"column(time_elapsed)"`
}

func (job *Job) RenderTemplate(ctx *context.Context) string {
	header, err := json.Marshal(ctx.Request.Header)
	query, err := json.Marshal(ctx.Input.Params())
	if err != nil {
		fmt.Println(err)
	}
	tmplData := templates.TemplateData{Code: job.Code, Header: string(header), Query: string(query)}
	tmpl, err := template.Must(template.New("php").Parse("<?php class Request {     public $headers;     public $query; } {{.Code}} $request = new Request(); $request->headers = json_decode('{{.Header}}'); $request->query = json_encode('{{.Query}}'); echo main($request);")).ParseFiles("templates/" + job.Image + ".tmpl")
	if err != nil {
		log.Fatalf("Fehler beim Laden des Templates: %v", err)
	}
	var result bytes.Buffer
	err = tmpl.Execute(&result, tmplData)
	filename := time.Now().UnixNano()
	os.WriteFile(strconv.FormatInt(filename, 10), []byte(html.UnescapeString(string(result.Bytes()))), 0777)
	return strconv.FormatInt(filename, 10)
}

func RunJob(ctx *context.Context) {

	name := ctx.Input.Param(":name")
	fmt.Println("attempting to run " + name)
	var job Job
	o := orm.NewOrm()
	o.QueryTable((*Job)(nil)).Filter("name", name).One(&job)

	if ctx.Request.Header.Get("token") == job.Token {
		startTime := time.Now()

		codefile := job.RenderTemplate(ctx)

		cmd := exec.Command("docker", []string{"run", "--rm", "-v", "./" + codefile + ":/code", job.Image, "php", "/code"}...)
		fmt.Println(cmd.Args)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(string(output))
			log.Fatalf("Befehlsausführung fehlgeschlagen: %v", err.Error())
		}
		os.Remove(codefile)
		ctx.WriteString(string(output))
		run := new(Run)
		run.Job = &job
		run.DateTime = time.Now().Format("2006-01-02 15:04:05")
		run.TimeElapsed = time.Since(startTime).Milliseconds()
		o.Insert(run)
		fmt.Println("done in " + strconv.FormatInt(time.Since(startTime).Milliseconds(), 10) + " ms")
	} else {
		fmt.Println("wrong token!")
	}
}
