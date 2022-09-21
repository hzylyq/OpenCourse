// Contains the implementation of a LSP client.

package lsp

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	
	"github.com/cmu440/lspnet"
)

const MaxLen = 1024

type client struct {
	// base
	conn      *lspnet.UDPConn
	connId    int
	nowSeqNum int
	param     *Params
	
	// send
	
	// recv
	noMsgEpochCnt int
	
	epochTimer *time.Timer
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
	rAddr, err := lspnet.ResolveUDPAddr("udp", hostport)
	if err != nil {
		return nil, err
	}
	
	conn, err := lspnet.DialUDP("udp", nil, rAddr)
	if err != nil {
		return nil, err
	}
	
	c := &client{
		conn:       conn,
		connId:     0,
		nowSeqNum:  1,
		param:      params,
		epochTimer: time.NewTimer(0),
	}
	
	readBytes := make([]byte, MaxLen)
	for {
		select {
		case <-c.epochTimer.C:
			c.noMsgEpochCnt++
			if c.noMsgEpochCnt >= c.param.EpochLimit {
				return nil, fmt.Errorf("Error Connection not established")
			}
			
			// todo send msg
			c.SendMsg()
			c.epochTimer.Reset(time.Millisecond * time.Duration(c.param.EpochMillis))
		default:
			msg, rAddr, err := c.RecvMsg(readBytes)
			if err != nil {
				log.Println(err)
				continue
			}
			
			c.noMsgEpochCnt = 0
			
		}
	}
	
	return c, nil
}

func (c *client) ConnID() int {
	return c.connId
}

func (c *client) Read() ([]byte, error) {
	reader := bufio.NewReader(&c.conn)
	
	for {
		reader.
		()
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

func (c *client) SendMsg(msg *Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	
	_, err = c.conn.Write(b)
	return err
}

func (c *client) RecvMsg(b []byte) (*Message, *lspnet.UDPAddr, error) {
	readSize, rAddr, err := c.conn.ReadFromUDP(b)
	if err != nil {
		return nil, nil, err
	}
	
	var msg Message
	err := json.Unmarshal(b[:readSize], &msg)
	if err != nil {
		return nil, nil, err
	}
	
	return &msg, rAddr, nil
}
