// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

type client struct {
	conn net.Conn

	s *multiEchoServer
}

const MaxQueue = 100

type multiEchoServer struct {
	listener net.Listener

	message chan string

	clientList []*client

	count int
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	return &multiEchoServer{
		message: make(chan string, 100),
	}
}

func (mes *multiEchoServer) Start(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	mes.listener = listener

	go mes.boardCast()

	go func() {
		for {
			conn, err := mes.listener.Accept()
			if err != nil {
				panic(err)
			}

			c := &client{
				conn: conn,
				s:    mes,
			}

			mes.clientList = append(mes.clientList, c)
			go c.Read()
		}
	}()

	go http.ListenAndServe(":10000", nil)

	return nil
}

func (mes *multiEchoServer) Close() {
	close(mes.message)
	mes.listener.Close()
}

func (mes *multiEchoServer) Count() int {
	return len(mes.clientList)
}

func (mes *multiEchoServer) boardCast() {
	for message := range mes.message {
		for _, cli := range mes.clientList {
			log.Println(message)
			cli.conn.Write([]byte(message))
		}
	}
}

func (c *client) Read() {
	reader := bufio.NewReader(c.conn)

	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			panic(err)
		}

		if len(c.s.message) > MaxQueue {
			return
		}

		c.s.message <- string(msg)
	}
}
