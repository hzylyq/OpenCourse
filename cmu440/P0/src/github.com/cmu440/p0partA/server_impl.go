// Implementation of a KeyValueServer. Students should write their code in this file.

package p0partA

import (
	"fmt"
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

func (kvs *keyValueServer) Start(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	for {
		listener.Accept()
	}

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
