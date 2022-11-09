package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

const (
	MessageRequest = 1
	FinishRequest  = 2
)

const (
	UnAssigned = 0
	Assigned   = 1
)

type MapReduceArgs struct {
	MessageType int
	Task        MapReduceTask
}

type MapReduceReply struct {
	Task MapReduceTask
}

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type Arg struct {
}

type Reply struct {
	Job
}

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

type FinishArg struct {
	Job
}

type FinishReply struct {
	Files []string
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
