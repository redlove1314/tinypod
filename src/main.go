package main

import (
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

const (
	NOT_FOUND = `
<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
<html><head>
<title>404 Not Found</title>
</head><body>
<h1>Not Found</h1>
<p>The requested URL #?# was not found on this server.</p>
</body></html>
`

	METHOD_NOT_ALLOWED = `
<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
<html><head>
<title>405 Method Not Allowed</title>
</head><body>
<h1>Method Not Allowed</h1>
<p>The #?# method is not allowed for the requested URL.</p>
</body></html>
`
)

var (
	port        = 80
	contextPath = ""
)

func main() {

	initClientFlags()

	workDir, _ := os.Getwd()
	if contextPath != "" {
		absPath, _ := filepath.Abs(contextPath)
		contextPath = lib.FixPath(absPath)
	} else {
		absPath, _ := filepath.Abs(workDir)
		contextPath = lib.FixPath(absPath)
	}
	log.Print("mapping context path:", contextPath)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodOptions {
			writer.Header().Add("Access-Control-Allow-Method", "GET,OPTIONS")
			writer.Header().Add("Access-Control-Allow-Origin", "*")
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte(strings.Replace(METHOD_NOT_ALLOWED, "#?#", request.Method, -1)))
			return
		}
		log.Print("get->["+request.URL.Path, "]")
		filename := contextPath + request.URL.Path
		// fmt.Println(request.RequestURI)
		if !lib.Exists(filename) {
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(strings.Replace(NOT_FOUND, "#?#", request.URL.Path, -1)))
			return
		}
		f, _ := lib.GetFile(filename)
		io.Copy(writer, f)
	})
	s := &http.Server{
		Addr: ":" + strconv.Itoa(port),
		// ReadTimeout:    10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      0,
		MaxHeaderBytes:    1 << 20,
	}
	log.Print("http server listening on port:", strconv.Itoa(port))
	err := s.ListenAndServe()
	if err != nil {
	}
	log.Fatal("error:", err)
}

func initClientFlags() {

	appFlag := cli.NewApp()
	appFlag.Name = "go httpd!"
	appFlag.Usage = ""

	// config file location
	appFlag.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "port,p",
			Value:       80,
			Usage:       "http port",
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "context,c",
			Value:       "",
			Usage:       "mapping context directory",
			Destination: &contextPath,
		},
	}

	appFlag.Action = func(c *cli.Context) error {
		return nil
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
