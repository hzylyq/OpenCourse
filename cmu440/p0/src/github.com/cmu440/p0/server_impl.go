// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"fmt"
	"log"
	"net"
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

	for {
		conn, err := mes.listener.Accept()
		if err != nil {
			return err
		}

		c := &client{
			conn: conn,
			s:    mes,
		}

		mes.clientList = append(mes.clientList, c)
		go c.Read()
	}

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
	b := make([]byte, 1024)
	if _, err := c.conn.Read(b); err != nil {
		log.Fatal(err)
	}

	log.Println(string(b))

	if len(c.s.message) > MaxQueue {
		return
	}

	c.s.message <- string(b)
}
