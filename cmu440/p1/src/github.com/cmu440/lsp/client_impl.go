// Contains the implementation of a LSP client.

package lsp

import (
	"bufio"
	"errors"
	"net"
)

type client struct {
	conn net.UDPConn
	id   int
}

// NewClient creates, initiates, and returns a new client. This function
// should return after a connection with the server has been established
// (i.e., the client has received an Ack message from the server in response
// to its connection request), and should return a non-nil error if a
// connection could not be made (i.e., if after K epochs, the client still
// hasn't received an Ack message from the server in response to its K
// connection requests).
//
// hostport is a colon-separated string identifying the server's host address
// and port number (i.e., "localhost:9999").
func NewClient(hostport string, params *Params) (Client, error) {
	return nil, errors.New("not yet implemented")
}

func (c *client) ConnID() int {
	return c.id
}

func (c *client) Read() ([]byte, error) {
	reader := bufio.NewReader(&c.conn)

	for {
		reader.()
	}

	// TODO: remove this line when you are ready to begin implementing this method.
	select {} // Blocks indefinitely.
	return nil, errors.New("not yet implemented")
}

func (c *client) Write(payload []byte) error {
	return errors.New("not yet implemented")
}

func (c *client) Close() error {
	c.conn.Close()

	return errors.New("not yet implemented")
}
