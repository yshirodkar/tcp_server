package main

import (
    "fmt"
    "net"
    "os"
    "errors"
    "bytes"
    // "encoding/json"
    "encoding/gob"
    // "bufio"
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
        //logs an incoming message
        fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())

        // Handle connections in a new goroutine.
        go handleRequest(conn)
        fmt.Println("I am done ")
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
    // Make a buffer to hold incoming data.
    netData := make([]byte, 129)
    res := make([]byte, 1)
    // Read the incoming connection into the buffer.
    lenVal, err := conn.Read(netData)
    // netData, err := bufio.NewReader(conn).ReadBytes('\n')
    fmt.Println("The input is ", netData)
    fmt.Println("The length is ", lenVal)
        if err != nil {
            fmt.Println("Error reading:", err.Error())
        }
        msb1 := (netData[0] & 0xff) >> 7;
        if msb1 == 0 {
            fmt.Println("Sending to push")
            Push(netData[:netData[0]+1])
            // Send a response back to person contacting us.
            conn.Write(res)
        } else if msb1 == 1 {
            fmt.Println("poping it", stack)
            popRes, err := Peek()
            fmt.Println("after poping it", stack)
            if err != nil {
                fmt.Println("Error while popping stack", err.Error())
            } else {
                b, ok := popRes.([]byte)
                if !ok {
                    fmt.Println("This is wrong")
                }
                if b != nil {
                    resPopb := b[:b[0]+1]
                    fmt.Println("the response is ", resPopb)
                    // Send a response back to person contacting us.
                    conn.Write(resPopb)
                }
                conn.Close()
            }
        }

    // Close the connection when you're done with it.
    conn.Close()
}

// Push ...
func Push(element interface{}) {
    fmt.Println("push element ", element)
    if len(stack) < 100 {
        stack = append(stack, element)        
    }
    fmt.Println("stack length", len(stack))
}

// Pop removes the last element of this stack. If stack is empty, it returns
// -1 and an error.
func  Pop() (interface{}, error) {
    fmt.Println("stack is ", stack)
    if len(stack) > 0 {
        popped := &(stack)[len(stack)-1]
        stack = (stack)[:len(stack)-1]
        fmt.Println("testing ", popped)
        return popped, nil
    }
    return -1, errors.New("stack is empty")
}

// Peek returns the topmost element of the stack. If stack is empty, it returns
// -1 and an error.
func Peek() (interface{}, error) {
    fmt.Println("stack is ", stack)
    if len(stack) > 0 {
        popped := (stack)[len(stack)-1]
        stack = stack[:len(stack)-1]
        return popped, nil
    }
 return -1, errors.New("stack is empty")
}

func MarshalBytes(i interface{}) ([]byte, error) {

    var buf bytes.Buffer

    enc := gob.NewEncoder(&buf)

    if err := enc.Encode(i); err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}
