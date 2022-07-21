// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

type client struct {
	conn net.Conn

	s *multiEchoServer

	Id int64

	recv chan string

	closeCh chan struct{}
}

const MaxQueue = 100

type multiEchoServer struct {
	listener net.Listener

	message chan string

	clients chan map[int64]*client

	close chan struct{}
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {

	return &multiEchoServer{
		message: make(chan string, 100),
		clients: make(chan map[int64]*client, 1),
		close:   make(chan struct{}, 1),
	}
}

func (mes *multiEchoServer) Start(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	clientMap := make(map[int64]*client)
	mes.clients <- clientMap

	mes.listener = listener

	go mes.boardCast()

	go func() {
		var miniId int64

		for {
			select {
			case <-mes.close:
				return
			default:
				break
			}

			conn, err := mes.listener.Accept()
			if err != nil {
				log.Println(err)
				return
			}

			c := &client{
				conn: conn,
				s:    mes,
				Id:   miniId,
			}

			clientMap, ok := <-mes.clients
			if !ok {
				clientMap = make(map[int64]*client)
			}

			clientMap[miniId] = c
			miniId++

			go func() {
				mes.clients <- clientMap
			}()

			go c.Read()
			// go c.Write()
		}
	}()

	go http.ListenAndServe(":10000", nil)

	return nil
}

func (mes *multiEchoServer) Close() {
	mes.close <- struct{}{}
	close(mes.message)
	mes.listener.Close()
}

func (mes *multiEchoServer) Count() int {
	cli := <-mes.clients
	count := len(cli)

	mes.clients <- cli

	return count
}

func (mes *multiEchoServer) boardCast() {
	for message := range mes.message {
		clients := <-mes.clients

		for _, cli := range clients {
			cli.conn.Write([]byte(message))
		}

		mes.clients <- clients
	}
}

func (c *client) Read() {
	reader := bufio.NewReader(c.conn)

	for {
		msg, err := reader.ReadBytes('\n')
		switch err {
		case nil:
			break
		case io.EOF:
			c.conn.Close()

			cliMap := <-c.s.clients

			delete(cliMap, c.Id)

			c.s.clients <- cliMap

			return
		default:
			panic(err)
		}

		// todo when close
		// select {
		// case <-c.s.close:
		// 	return
		// }

		if len(c.s.message) > MaxQueue {
			return
		}

		c.s.message <- string(msg)
	}
}

func (c *client) Write() {
	select {
	case msg := <-c.recv:
		c.conn.Write([]byte(msg))
	case <-c.closeCh:
		return
	}
}
