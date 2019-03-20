package main

import (
	"common"
	"github.com/urfave/cli"
	"io"
	"log"
	"net"
	"os"
	"regexp"
)

func main() {

	initProxyFlags()

	if common.Command != common.COMMAND_START {
		return
	}

	addrPattern := "^(.*:)([0-9]+)$"

	if mat, err := regexp.Match(addrPattern, []byte(common.LocalAddr)); err != nil && mat {
		// legal
	} else {
		log.Fatal("error: illegal local address")
	}
	if mat, err := regexp.Match(addrPattern, []byte(common.RemoteAddr)); err != nil && mat {
		// legal
	} else {
		log.Fatal("error: illegal remote address")
	}

	listener, err := net.Listen("tcp", common.LocalAddr)
	if err != nil {
		log.Fatal("error listening on local port: ", err)
	}
	log.Println("server listening at ", common.LocalAddr)
	for {
		conn, err := listener.Accept()
		// log.Println("new connection...")
		if err != nil {
			log.Fatal("error accept connection: ", err)
		}
		backendConn, err := net.Dial("tcp", common.RemoteAddr)
		if err != nil {
			log.Fatal("error connect to server: ", err)
		}
		pipe(backendConn, conn)
	}
}

func pipe(conn1 net.Conn, conn2 net.Conn) {
	go func() {
		_, err := io.Copy(conn1, conn2)
		if err != nil {
			// log.Println("error: ", err, ", read bytes:", len)
		}
		// fmt.Println("pipe end1")
		conn1.Close()
		conn2.Close()
	}()
	go func() {
		_, err := io.Copy(conn2, conn1)
		if err != nil {
			// log.Println("error: ", err, ", read bytes:", len)
		}
		// fmt.Println("pipe end2")
		conn1.Close()
		conn2.Close()
	}()
}

func initProxyFlags() {

	appFlag := cli.NewApp()
	appFlag.Name = "go proxy!"
	appFlag.Usage = ""
	appFlag.Version = "1.0"

	// sub command 'upload'
	appFlag.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "map a local port to a remote port",
			Action: func(c *cli.Context) error {
				common.Command = common.COMMAND_START
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "local,l",
					Value:       "",
					Usage:       "listen address(host:port)",
					Destination: &common.LocalAddr,
				},
				cli.StringFlag{
					Name:        "remote,r",
					Value:       "",
					Usage:       "remote address(host:port)",
					Destination: &common.RemoteAddr,
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
