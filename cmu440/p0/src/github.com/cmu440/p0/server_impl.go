// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"fmt"
	"log"
	"net"
)

type multiEchoServer struct {
	listener net.Listener

	message chan []byte

	connList []net.Conn

	count int
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	return &multiEchoServer{
		message: make(chan []byte, 100),
	}
}

func (mes *multiEchoServer) Start(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	mes.listener = listener

	for {
		conn, err := mes.listener.Accept()
		if err != nil {
			return err
		}

		mes.connList = append(mes.connList, conn)

		go func() {
			var b []byte
			if _, err = conn.Read(b); err != nil {
				log.Fatal(err)
			}

			mes.message <- b

			// mes.count +=

		}()
	}

	mes.boardCast()

	return nil
}

func (mes *multiEchoServer) Close() {
	close(mes.message)
	mes.listener.Close()
}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!
	return -1
}

func (mes *multiEchoServer) boardCast() {
	for message := range mes.message {
		for _, conn := range mes.connList {
			conn.Write(message)
		}
	}
}
