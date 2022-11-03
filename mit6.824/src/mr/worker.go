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
	"time"
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
			reduceTask(reducef, reply.Task)
		case WaitJob:
			waitTask()
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

	kvaa := make([][]KeyValue, task.NumReduce)

	for _, kv := range kva {
		idx := ihash(kv.Key) % task.NumReduce
		kvaa[idx] = append(kvaa[idx], kv)
	}

	for i := range kvaa {
		storeIntermediateFile(kvaa[i], intermediateFileName(task.TaskNum, i))
	}

	finishTask(task)
}

func reduceTask(reduceF func(string, []string) string, task MapReduceTask) {
	log.Printf("task is reduce job")

	var intermediate []KeyValue
	for _, filename := range task.ReduceFiles {
		intermediate = append(intermediate, loadIntermediateFile(filename)...)
	}

	sort.Sort(ByKey(intermediate))

	oname := fmt.Sprintf("mr-out-%d", task.TaskNum)
	ofile, _ := os.Create(oname)

	//
	// call Reduce on each distinct key in intermediate[],
	// and print the result to mr-out-0.
	//
	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := reduceF(intermediate[i].Key, values)

		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}

	ofile.Close()

	finishTask(task)
}

func waitTask() {
	time.Sleep(1 * time.Second)
}

func storeIntermediateFile(kva []KeyValue, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("can open %v", fileName)
	}
	defer file.Close()

	enc := json.NewEncoder(file)

	for _, kv := range kva {
		err = enc.Encode(&kv)
		if err != nil {
			log.Fatal("encode failed")
		}
	}
}

func intermediateFileName(numMapTask, numReduceTask int) string {
	return fmt.Sprintf("mr-%d-%d", numMapTask, numReduceTask)
}

func loadIntermediateFile(fileName string) []KeyValue {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("cannot open %v", fileName)
	}
	defer file.Close()

	kva := make([]KeyValue, 0)
	dec := json.NewDecoder(file)
	for {
		kv := KeyValue{}

		if err = dec.Decode(&kv); err != nil {
			log.Fatal("decode failed")
		}
		kva = append(kva, kv)
	}

	return kva
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

func finishTask(task MapReduceTask) {
	args := MapReduceArgs{
		MessageType: FinishRequest,
		Task:        task,
	}
	reply := MapReduceReply{}
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		fmt.Printf("finished")
	} else {
		fmt.Printf("call failed!")
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
