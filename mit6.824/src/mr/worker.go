package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"sort"
)

type Work struct {
	Idx   int
	State int
}

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

type ByKey []KeyValue

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// d
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// Worker run mapReduce
// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	//
	// a big difference from real MapReduce is that all the
	// intermediate data is in one place, intermediate[],
	// rather than being partitioned into NxM buckets.
	//

	for {
		args := MapReduceArgs{
			MessageType: MessageRequest,
			Task:        MapReduceTask{},
		}

		reply := MapReduceReply{}

		ok := call("Coordinator.Task", &args, &reply)
		if !ok {
			log.Println("call rpc failed")
			break
		}
		log.Printf("%v", reply)

		switch reply.Task.TaskType {
		case MapJob:
			mapTask(mapf, reply.Task)

		case ReduceJob:
			log.Printf("task is reduce job")

			// todo open files

			var intermediate []KeyValue
			for _, file := range reply.File {
				f, err := os.Open(file)
				if err != nil {
					panic(err)
				}

				dec := json.NewDecoder(f)
				for {
					var kv KeyValue
					if err = dec.Decode(&kv); err != nil {
						break
					}

					intermediate = append(intermediate, kv)
				}

				f.Close()
			}

			oFile, err := os.CreateTemp("./", fmt.Sprintf("mr-out-%d", reply.Id))
			if err != nil {
				panic(err)
			}

			i := 0
			for i < len(intermediate) {
				j := i + 1
				for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
					j++
				}
				var values []string
				for k := i; k < j; k++ {
					values = append(values, intermediate[k].Value)
				}
				output := reducef(intermediate[i].Key, values)

				// this is the correct format for each line of Reduce output.
				fmt.Fprintf(oFile, "%v %v\n", intermediate[i].Key, output)

				i = j
			}

			oFile.Close()
		case WaitJob:

		}

		if reply.JobType == MapJob || reply.JobType == ReduceJob {
			Finish(reply.Job)
		} else {
			break
		}

	}

	log.Println("all job finished")

	// uncomment to send the Example RPC to the coordinator.
	// var reply mapTaskReply
	// call("Coordinator.Map", &mapTaskArg{}, &reply)

	// keyValues := mapf(reply.File, reply.)

}

func mapTask(mapF func(string, string) []KeyValue, task MapReduceTask) {
	log.Printf("task is map job")
	intermediate := make([]KeyValue, 0)

	fileName := task.MapFile

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("cannot open %v", fileName)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", fileName)
	}
	file.Close()
	kva := mapF(fileName, string(content))
	intermediate = append(intermediate, kva...)

	for _, filename := range reply.File {

	}
	sort.Sort(ByKey(intermediate))

	nReduce := reply.nReduce

	outPutFiles := make([]*os.File, reply.nReduce)
	fileEncs := make([]*json.Encoder, nReduce)

	for i := 0; i < nReduce; i++ {
		outPutFiles[i], _ = os.CreateTemp("mr-tmp", fmt.Sprintf("mr-%d-%d", reply.Id, i))
		fileEncs[i] = json.NewEncoder(outPutFiles[i])
	}

	for _, kv := range intermediate {
		i := ihash(kv.Key) % reply.nReduce
		enc := fileEncs[i]

		err := enc.Encode(&kv)
		if err != nil {
			panic("json encode failed")
		}
	}
}

// CallExample example function to show how to make an RPC call to the coordinator.
// the RPC argument and reply types are defined in rpc.go.
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

func Finish(job Job) {
	args := FinishArg{Job: job}
	reply := FinishReply{}
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		fmt.Printf("finished")
	} else {
		fmt.Printf("call failed!\n")
	}
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
