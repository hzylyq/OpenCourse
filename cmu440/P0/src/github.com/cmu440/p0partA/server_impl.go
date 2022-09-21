// Implementation of a KeyValueServer. Students should write their code in this file.

package p0partA

import (
	"fmt"
	"log"
	"net"
	
	"github.com/cmu440/p0partA/kvstore"
)

type keyValueServer struct {
	store kvstore.KVStore
	
	ln net.Listener
}

// New creates and returns (but does not start) a new KeyValueServer.
func New(store kvstore.KVStore) KeyValueServer {
	s := &keyValueServer{
		store: store,
	}
	
	return s
}

type client struct {
	id   int64
	conn net.Conn
}

func (kvs *keyValueServer) Start(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				return
			}
			
			go
			
		}
		
	}()
	
	return nil
}

func (kvs *keyValueServer) Close() {
	// TODO: implement this!
}

func (kvs *keyValueServer) CountActive() int {
	// TODO: implement this!
	return -1
}

func (kvs *keyValueServer) CountDropped() int {
	// TODO: implement this!
	return -1
}

// TODO: add additional methods/functions below!
