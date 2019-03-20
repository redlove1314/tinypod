package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal("error listening on local port: ", err)
	}

	for {
		conn, err := listener.Accept()
		log.Println("new connection...")
		if err != nil {
			log.Fatal("error accept connection: ", err)
		}
		backendConn, err := net.Dial("tcp", "foxmaz.com:27754")
		if err != nil {
			log.Fatal("error connect to server: ", err)
		}
		Pipe(backendConn, conn)
	}
}

func chanFromConn(conn net.Conn) chan []byte {
	c := make(chan []byte)

	go func() {
		b := make([]byte, 1024)
		for {
			n, err := conn.Read(b)
			if n > 0 {
				res := make([]byte, n)
				copy(res, b[:n])
				c <- res
			}
			if err != nil {
				log.Println("error: ", err)
				c <- nil
				break
			}
		}
	}()
	return c
}

func Pipe(conn1 net.Conn, conn2 net.Conn) {
	chan1 := chanFromConn(conn1)
	chan2 := chanFromConn(conn2)

	for {
		select {
		case b1 := <-chan1:
			if b1 == nil {
				log.Println("connection close")
				return
			} else {
				conn2.Write(b1)
			}
		case b2 := <-chan2:
			if b2 == nil {
				log.Println("connection close")
				return
			} else {
				conn1.Write(b2)
			}

		}
	}
}
