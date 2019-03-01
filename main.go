package main

import (
	"errors"
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
	CONN_TYPE = "tcp"
)

type Stack []interface{}

var stack Stack

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	netData := make([]byte, 128)
	res := make([]byte, 1)
	// Read the incoming connection into the buffer.
	lenVal, err := conn.Read(netData)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
		conn.Write(nil)
	} else {
		if lenVal < 1 {
			conn.Write(res)
		} else {
			msb1 := (netData[0] & 0xff) >> 7
			if msb1 == 0 {
				Push(netData[:netData[0]+1])
				// Send a response back to person contacting us.
				conn.Write(res)
			} else if msb1 == 1 {
				if len(stack) > 0 {
					popRes, err := Pop()
					if err != nil {
						fmt.Println("Error while popping stack", err.Error())
					} else {
						b, _ := popRes.([]byte)
						if b != nil {
							resPopb := b[:b[0]+1]
							// Send a response back to person contacting us.
							conn.Write(resPopb)
						}
					}
				} else {
					conn.Write(res)
				}
			}
			netData = make([]byte, 128)
		}
	}
	// Close the connection when you're done with it.
	conn.Close()
}

// Push ...
func Push(element interface{}) {
	if len(stack) < 100 {
		stack = append(stack, element)
	}
}

// Pop removes the last element of this stack. If stack is empty, it returns
// -1 and an error.
func Pop() (interface{}, error) {
	if len(stack) > 0 {
		popped := (stack)[len(stack)-1]
		stack = stack[:len(stack)-1]
		return popped, nil
	}
	return -1, errors.New("stack is empty")
}
