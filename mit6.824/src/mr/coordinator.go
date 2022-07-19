package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

const (
	Idle      = 1
	InProcess = 2
	Complete  = 3
)

const (
	MapJob    = 1
	ReduceJob = 2
	FinishJob = 3
)

// Job map or reduce job
type Job struct {
	JobType int
	File    []string
	Id      int64
	State   int
	nReduce int
}

type Coordinator struct {
	// Your definitions here.
	// master 节点
	Idx      int
	IsFinish bool

	jobTask   chan Job
	reduceNum int
	nReduce   int

	tmpFiles []string

	wg sync.WaitGroup
}

// Your code here -- RPC handlers for the worker to call.
func (c *Coordinator) GenMapTask() {

}

func (c *Coordinator) Task(arg *Arg, reply *Reply) error {
	job := <-c.jobTask
	reply.Job = job

	go func() {
		job.State = InProcess
		c.jobTask <- job
	}()
	// todo timeout check

	return nil
}

func (c *Coordinator) Finish(arg FinishArg, reply FinishReply) error {
	switch arg.JobType {
	case MapJob:
		if len(c.jobTask) == 0 {
			// generate reduce task
		}

		c.tmpFiles = append(c.tmpFiles, reply.Files...)

	case ReduceJob:

	}

	return nil
}

func (c *Coordinator) MakeReduce() {
	reduceFileMap := make(map[int64][]string)

	for id, tmpFile := range reduceFileMap {
		job := Job{
			Id:    id,
			State: Idle,
			File:  tmpFile,
		}

		c.jobTask <- job
	}
}

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	// l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// Done
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	// ret := false

	// Your code here.

	return c.IsFinish
}

// MakeCoordinator create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	// Your code here.

	c := Coordinator{
		nReduce: nReduce,
	}

	// generate task queue
	taskList := make([]Job, len(files))

	for i, file := range files {
		task := taskList[i]

		task.Id = int64(i)
		task.JobType = MapJob
		task.State = Idle
		task.File = append(task.File, file)
		task.nReduce = nReduce

		taskList[i] = task
	}

	log.Printf("%v", taskList)

	c.jobTask = make(chan Job, nReduce)
	for _, task := range taskList {
		c.jobTask <- task
	}

	c.server()
	return &c
}

func (c *Coordinator) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
