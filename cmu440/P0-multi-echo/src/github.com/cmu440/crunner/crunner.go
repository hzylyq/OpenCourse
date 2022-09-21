package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	defaultHost = "localhost"
	defaultPort = 9999
)

// To test your server implementation, you might find it helpful to implement a
// simple 'client runner' program. The program could be very simple, as long as
// it is able to connect with and send messages to your server and is able to
// read and print out the server's echoed response to standard output. Whether or
// not you add any code to this file will not affect your grade.
func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", defaultPort))
	if err != nil {
		log.Fatalln(err)
	}

	buffer := make([]byte, 1024)

	// 准备命令行标准输入
	reader := bufio.NewReader(os.Stdin)

	for {
		lineBytes, _, _ := reader.ReadLine()
		conn.Write(lineBytes)
		n, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}

		serverMsg := string(buffer[0:n])
		fmt.Printf("服务端msg", serverMsg)
		if serverMsg == "bye" {
			break
		}
	}

}
