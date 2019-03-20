// httpd.go is a simple http server which can map a local directory as a http server.
// Usage:
// httpd -p 8080 -c /my/path
// httpd --port 8080 --context /my/path

package main

import (
	"common"
	"github.com/urfave/cli"
	"io"
	"lib"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {

	initHttpFlags()

	if common.Command != common.COMMAND_START {
		return
	}

	workDir, _ := os.Getwd()
	if common.ContextPath != "" {
		absPath, _ := filepath.Abs(common.ContextPath)
		common.ContextPath = lib.FixPath(absPath)
	} else {
		absPath, _ := filepath.Abs(workDir)
		common.ContextPath = lib.FixPath(absPath)
	}

	log.Print("mapping context path:", common.ContextPath)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		user, pass, _ := request.BasicAuth()
		if common.BasicAuth != "" && common.BasicAuth != user+":"+pass {
			writer.Header().Add("WWW-Authenticate", "Basic realm=\"example\"")
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(common.FORBIDDEN))
			return
		}

		if request.Method == http.MethodOptions {
			writer.Header().Add("Access-Control-Allow-Method", "GET,OPTIONS")
			writer.Header().Add("Access-Control-Allow-Origin", "*")
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte(strings.Replace(common.METHOD_NOT_ALLOWED, "#?#", request.Method, -1)))
			return
		}
		log.Print("get->["+request.URL.Path, "]")
		filename := common.ContextPath + request.URL.Path
		// fmt.Println(request.RequestURI)
		if !lib.Exists(filename) {
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(strings.Replace(common.NOT_FOUND, "#?#", request.URL.Path, -1)))
			return
		}
		f, _ := lib.GetFile(filename)
		io.Copy(writer, f)
	})

	s := &http.Server{
		Addr: ":" + strconv.Itoa(common.Port),
		// ReadTimeout:    10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      0,
		MaxHeaderBytes:    1 << 20,
	}
	log.Print("http server listening on port:", strconv.Itoa(common.Port))
	err := s.ListenAndServe()
	if err != nil {
	}
	log.Fatal("error:", err)
}

func initHttpFlags() {

	appFlag := cli.NewApp()
	appFlag.Name = "go http!"
	appFlag.Usage = ""
	appFlag.Version = "1.0"

	// sub command 'upload'
	appFlag.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start an http server",
			Action: func(c *cli.Context) error {
				common.Command = common.COMMAND_START
				return nil
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "port,p",
					Value:       80,
					Usage:       "http port",
					Destination: &common.Port,
				},
				cli.StringFlag{
					Name:        "context,c",
					Value:       "",
					Usage:       "mapping context directory",
					Destination: &common.ContextPath,
				},
				cli.StringFlag{
					Name:        "auth,a",
					Value:       "",
					Usage:       "http basic auth(such as \"admin:123456\")",
					Destination: &common.BasicAuth,
				},
			},
		},
	}

	// 帮助文件模板
	cli.AppHelpTemplate = `name:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}
usage:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .VisibleCommands}}
commands:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
options:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}
copyright:
   {{.Copyright}}{{end}}
`
	err := appFlag.Run(os.Args)
	if err != nil {
		log.Fatal("err:", err)
	}
}
