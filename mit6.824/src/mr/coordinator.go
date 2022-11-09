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
	WaitJob   = 3
)

type MapReduceTask struct {
	TaskType   int
	TaskStatus int
	TaskNum    int

	MapFile     string
	ReduceFiles []string

	NumReduce int
	NumMap    int
}

type Coordinator struct {
	NumMap            int
	NumMapFinished    int
	NumReduce         int
	NumReduceFinished int

	MapTasks    []MapReduceTask
	ReduceTasks []MapReduceTask

	MapFinish    bool
	ReduceFinish bool

	mu sync.Mutex
}

// GenMapTask generate task
// Your code here -- RPC handlers for the worker to call.
func (c *Coordinator) GenMapTask() {

}

func (c *Coordinator) MakeTask(args *MapReduceArgs, reply *MapReduceReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if args.MessageType == MessageRequest {
		if !c.MapFinish {
			for i, task := range c.MapTasks {
				if task.TaskStatus == UnAssigned {
					c.MapTasks[i].TaskStatus = Assigned
					reply.Task = c.MapTasks[i]

					return nil
				}
			}
			reply.Task.TaskType = WaitJob
			return nil
		}

		if !c.ReduceFinish {
			for i, task := range c.ReduceTasks {
				if task.TaskStatus == UnAssigned {
					c.ReduceTasks[i].TaskStatus = Assigned
					reply.Task = c.ReduceTasks[i]

					return nil
				}
			}

			reply.Task.TaskType = WaitJob
			return nil
		}
	}

	if args.MessageType == FinishRequest {
		if args.Task.TaskType == MapJob

	}

	return nil

}

// TimeoutCheck when 10s end, job should be finished, if not put it to job task
func (c *Coordinator) TimeoutCheck(j Job) {

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

	return c.State == FinishJob
}

// MakeCoordinator create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {

	c := Coordinator{
		nReduce: nReduce,
		State:   MapJob,
	}

	c.jobTask = make(chan Job, nReduce)

	// generate task queue
	taskList := make([]Job, len(files))

	for i, file := range files {
		task := taskList[i]

		task.Id = int64(i)
		task.JobType = MapJob
		task.State = Idle
		task.File = append(task.File, file)
		task.nReduce = nReduce

		c.jobTask <- task
	}

	c.server()
	return &c
}

func (c *Coordinator) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
