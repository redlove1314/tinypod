// httpd.go is a simple http server which can map a local directory as a http server.
// Usage:
// httpd -p 8080 -c /my/path
// httpd --port 8080 --context /my/path

package main

import (
	"bytes"
	"common"
	"container/list"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"lib"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var httpClient *http.Client

func main() {

	runtime.GOMAXPROCS(500)

	initHttpFlags()

	if common.Command != common.COMMAND_START {
		return
	}
	// parse backend mappings
	common.ParseBackendMapping()

	workDir, _ := os.Getwd()
	if common.WorkDir != "" {
		absPath, _ := filepath.Abs(common.WorkDir)
		common.WorkDir = lib.FixPath(absPath)
	} else {
		absPath, _ := filepath.Abs(workDir)
		common.WorkDir = lib.FixPath(absPath)
	}

	common.ContextPath = lib.FixPath("/" + common.ContextPath)
	log.Println("map local directory", common.WorkDir, "to http context", common.ContextPath)

	httpClient = &http.Client{}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if common.ContextPath != "/" {
			if len(request.URL.Path) < len(common.ContextPath) ||
				(len(request.URL.Path) == len(common.ContextPath) && request.URL.Path != common.ContextPath) ||
				(len(request.URL.Path) > len(common.ContextPath) &&
					(request.URL.Path[0:len(common.ContextPath)] != common.ContextPath || request.URL.Path[len(common.ContextPath):len(common.ContextPath)+1] != "/")) {
				lib.WriteHttpResponse(writer, http.StatusNotFound, strings.Replace(common.NOT_FOUND, "#?#", request.URL.Path, -1), nil)
				return
			}
		}
		user, pass, _ := request.BasicAuth()
		if common.BasicAuth != "" && common.BasicAuth != user+":"+pass {
			lib.WriteHttpResponse(writer, http.StatusUnauthorized, common.FORBIDDEN, map[string]string{"WWW-Authenticate": "Basic realm=\"example\""})
			return
		}

		if request.Method == http.MethodOptions {
			lib.WriteHttpResponse(writer, http.StatusNoContent, "",
				map[string]string{"Access-Control-Allow-Method": "GET,OPTIONS", "Access-Control-Allow-Origin": "*"})
			return
		}

		/*if request.Method != http.MethodGet {
			lib.WriteHttpResponse(writer, http.StatusMethodNotAllowed, strings.Replace(common.METHOD_NOT_ALLOWED, "#?#", request.Method, -1), nil)
			return
		}*/

		key := invokeBackend(request.URL.Path)
		if key != "" {
			resolveBackend(key, request.URL.Path, writer, request)
			return
		}
		statusCode := http.StatusOK
		defer func() {
			log.Print("[" + strconv.Itoa(statusCode) + "] -> " + request.URL.Path)
		}()

		prefix := common.ContextPath
		if prefix == "/" {
			prefix = ""
		}
		if request.URL.Path == prefix+"/icon/folder.png" {
			writer.Header().Add("Content-Type", "image/png")
			writer.Header().Add("Content-Length", strconv.Itoa(len(common.IconFolder)))
			writer.Write([]byte(common.IconFolder))
			return
		} else if request.URL.Path == prefix+"/icon/file.png" {
			writer.Header().Add("Content-Type", "image/png")
			writer.Header().Add("Content-Length", strconv.Itoa(len(common.IconFile)))
			writer.Write([]byte(common.IconFile))
			return
		}
		tmp1 := lib.FixPath(strings.Replace(request.URL.Path, common.ContextPath, "", 1))
		filename := common.WorkDir + "/" + tmp1
		filename = lib.FixPath(filename)
		// fmt.Println(request.RequestURI)
		if !lib.Exists(filename) {
			statusCode = http.StatusNotFound
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(strings.Replace(common.NOT_FOUND, "#?#", request.URL.Path, -1)))
			return
		}
		if lib.IsDir(filename) {
			if !common.IndexDir {
				statusCode = http.StatusForbidden
				writer.WriteHeader(http.StatusForbidden)
				writer.Write([]byte(common.FORBIDDEN))
				return
			}
			cnt := indexDir(lib.FixPath(request.URL.Path), filename)
			ret := strings.Replace(common.DIR_INDEX, "#TITLE#", request.URL.Path, -1)
			ret = strings.Replace(ret, "#HEAD#", "Directory Index: "+request.URL.Path, -1)
			ret = strings.Replace(ret, "#CONTENT#", cnt, -1)
			writer.Write([]byte(ret))
		} else {
			f, _ := lib.GetFile(filename)
			io.Copy(writer, f)
		}
	})

	s := &http.Server{
		Addr: ":" + strconv.Itoa(common.Port),
		// ReadTimeout:    10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      0,
		MaxHeaderBytes:    1 << 20,
	}

	log.Println("http server listening on port", strconv.Itoa(common.Port))
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
					Value:       "/",
					Usage:       "mapping context path",
					Destination: &common.ContextPath,
				},
				cli.StringFlag{
					Name:        "workdir,w",
					Value:       "",
					Usage:       "mapping local directory",
					Destination: &common.WorkDir,
				},
				cli.StringFlag{
					Name:        "auth,a",
					Value:       "",
					Usage:       "http basic auth(such as \"admin:123456\")",
					Destination: &common.BasicAuth,
				},
				cli.BoolFlag{
					Name:        "index,i",
					Usage:       "whether support directory indexing",
					Destination: &common.IndexDir,
				},
				cli.StringFlag{
					Name:  "backend,b",
					Value: "",
					Usage: `proxy backend server api.
	for example, if you want to replace api 
	'http://xxx.com/api/anything'
	with 
	'http://localhost:8080/api/anything', you can do like this: 
	-b /api:http://xxx.com/api
		
	if you want to replace api
	'http://xxx.com/api/anything'
	with 
	'http://localhost:8080/3rd/api/anything', you can do like this:
	-b /3rd/api:http://xxx.com/api
	if you have many mappings, split with semicolon.
`,
					Destination: &common.BackendSettings,
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

func indexDir(uri string, path string) string {
	var buff bytes.Buffer
	buff.WriteString("<table><tr><th width='30px' align='center'></th><th align='left'>File Name</th><th>Size</th><th>Last Modified</th></tr>")
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return "Error index directory: " + path
	}

	var dirs list.List
	var files list.List
	for _, info := range infos {
		if info.IsDir() {
			dirs.PushBack(info)
		} else {
			files.PushBack(info)
		}
	}
	dirs.PushBackList(&files)

	for ele := dirs.Front(); ele != nil; ele = ele.Next() {
		info := ele.Value.(os.FileInfo)
		buff.WriteString("<tr>")
		prefix := common.ContextPath
		if prefix == "/" {
			prefix = ""
		}
		if info.IsDir() {
			buff.WriteString("<td><img class='icon' src='" + prefix + "/icon/folder.png'></td>")
		} else {
			buff.WriteString("<td><img class='icon' src='" + prefix + "/icon/file.png'></td>")
		}
		buff.WriteString("<td><a href=\"")
		buff.WriteString(uri)
		if uri != "/" {
			buff.WriteString("/")
		}
		buff.WriteString(info.Name())
		buff.WriteString("\">")
		buff.WriteString(info.Name())
		buff.WriteString("</a></td><td class='size'>")
		if info.IsDir() {
			buff.WriteString("-")
		} else {
			buff.WriteString(humanReadableSize(info.Size()))
		}
		buff.WriteString("</td><td class='time'>")
		buff.WriteString(lib.GetLongDateString(info.ModTime()))
		buff.WriteString("</td><td>")
		buff.WriteString("</tr>")
	}

	buff.WriteString("</table>")
	return buff.String()
}

func humanReadableSize(size int64) string {
	if size < 1024 {
		return strconv.FormatInt(size, 10) + "b"
	} else if size < 1048576 {
		return strconv.FormatFloat(float64(size)/1024, 'f', 2, 64) + "Kb"
	} else if size < 1073741824 {
		return strconv.FormatFloat(float64(size)/1048576, 'f', 2, 64) + "Mb"
	} else if size < 1099511627776 {
		return strconv.FormatFloat(float64(size)/1073741824, 'f', 2, 64) + "Gb"
	} else {
		return strconv.FormatFloat(float64(size)/1099511627776, 'f', 2, 64) + "Tb"
	}
}

func resolveBackend(k string, uri string, writer http.ResponseWriter, request *http.Request) {
	v := common.BackendMappings[k]
	proxyUrl := ""
	if strings.HasSuffix(v, "/") {
		pick := strings.TrimSpace(uri[len(k):])
		if strings.HasPrefix(pick, "/") {
			proxyUrl = v + pick[1:]
		} else {
			proxyUrl = v + pick
		}
	} else {
		proxyUrl = v + "/" + uri[1:]
	}
	queryS := lib.TValue(request.URL.Query().Encode() == "", "", "?"+request.URL.Query().Encode()).(string)

	statusCode := http.StatusOK
	defer func() {
		log.Println("["+strconv.Itoa(statusCode)+"] proxy", proxyUrl+queryS)
	}()

	req, err := http.NewRequest(request.Method, proxyUrl+queryS, request.Body)
	if err != nil {
		statusCode = http.StatusInternalServerError
		lib.WriteHttpResponse(writer, http.StatusInternalServerError, common.INTERNAL_SERVER_ERROR, nil)
		return
	}
	req.Header = request.Header
	resp, err := httpClient.Do(req)
	if err != nil {
		statusCode = http.StatusInternalServerError
		lib.WriteHttpResponse(writer, http.StatusInternalServerError, common.INTERNAL_SERVER_ERROR, nil)
		return
	}
	if resp.Header != nil {
		for k, v := range resp.Header {
			writer.Header().Add(k, v[0])
		}
	}
	statusCode = resp.StatusCode
	writer.WriteHeader(resp.StatusCode)

	defer resp.Body.Close()
	if req.Method != http.MethodHead {
		lib.Try(func() {
			buff := make([]byte, 1024)
			for {
				len, err := resp.Body.Read(buff)
				if len > 0 {
					writer.Write(buff[0:len])
					continue
				}
				if err != nil {
					// log.Println(err)
					break
				}
			}
			// log.Println("body read finish")
		}, func(i interface{}) {
			log.Println("error proxy url: ", proxyUrl, ":", i)
		})
	}
}

func invokeBackend(uri string) string {
	if len(common.BackendMappings) == 0 {
		return ""
	}
	for k := range common.BackendMappings {
		if k == "/" {
			return k
		}
		if strings.HasPrefix(uri, k) {
			if len(uri) == len(k) && uri != k {
				continue
			}
			if len(uri) > len(k) && (uri[0:len(k)] != k || uri[len(k):len(k)+1] != "/") {
				continue
			}
			return k
		}
	}
	return ""
}
