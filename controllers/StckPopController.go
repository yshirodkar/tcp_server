package controllers

import (
	"errors"
	"net/http"
	"philo_server/common"
)

// Stack is the Go implementation of stack
type Stack []interface{}

type IStackPopController interface {
	Push(element interface{})
	Pop() (interface{}, error)
	Peek() (interface{}, error)
}

type stackPopController struct {
	controller
}

func GetStackPopController() IStackPopController {
	return &stackPopController{}
}

// Push ...
func (s *Stack) Push(element interface{}) {
	if len(*s) < 100 {
		*s = append(*s, element)
	}
}

// Pop removes the last element of this stack. If stack is empty, it returns
// -1 and an error.
func (s *Stack) Pop() (interface{}, error) {
	if len(*s) > 0 {
		popped := (*s)[len(*s)-1]
		*s = (*s)[:len(*s)-1]
		return popped, nil
	}
	return -1, errors.New("stack is empty")
}

// Peek returns the topmost element of the stack. If stack is empty, it returns
// -1 and an error.
func (s *Stack) Peek() (interface{}, error) {
	if len(*s) > 0 {
		return (*s)[len(*s)-1], nil
	}
	return -1, errors.New("stack is empty")
}
