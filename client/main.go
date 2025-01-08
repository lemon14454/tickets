package main

import (
	"flag"
	"fmt"
)

var serverIP string
var serverPort int

// command line parser
// ./client -ip 127.0.0.1 -port	8080

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "server ip [default:127.0.0.1]")
	flag.IntVar(&serverPort, "port", 8080, "server port [default:8080]")
}

func main() {
	flag.Parse()
	client := NewClient(serverIP, serverPort)
	events, err := client.Request.GetAllEvents()
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fmt.Printf("%v", events)
}
