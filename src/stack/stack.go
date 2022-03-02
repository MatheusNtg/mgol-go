package stack

import "fmt"

var (
	ErrNotEnoughStackSpace = fmt.Errorf("no space left on the stack")
	ErrEmptyStack          = fmt.Errorf("empty stack")
)

type Stack struct {
	data     []interface{}
	capacity int
	length   int
}

//NewStack returns a new empty stack with
//maximum capacity defined by capacity
func NewStack(capacity int) *Stack {
	return &Stack{
		data:     nil,
		capacity: capacity,
		length:   0,
	}
}

//Push pushes the element present on the
//parameter element onto the stack. This
//is only possible if there is enough space
//on the stack, otherwise it will return an
//error and no element will be pushed on the
//stack
func (s *Stack) Push(element interface{}) error {
	if s.length >= s.capacity {
		return ErrNotEnoughStackSpace
	}

	s.data = append(s.data, element)
	s.length += 1
	return nil
}

//Pop returns and remove the top element of the stack
//This is only possible if the stack is not empty,
//otherwise it will return an error
func (s *Stack) Pop() (interface{}, error) {
	if s.length == 0 {
		return nil, ErrEmptyStack
	}
	element := s.data[s.length-1]
	s.data = s.data[:s.length-1]
	s.length -= 1
	return element, nil
}

//Get returns the top element of the stack
//This is only possible if the stack is not empty,
//otherwise it will return an error
func (s *Stack) Get() (interface{}, error) {
	if s.length == 0 {
		return nil, ErrEmptyStack
	}
	element := s.data[s.length-1]
	return element, nil
}

//Returns the length of stack this represents
//how many elements are present on the stack
func (s *Stack) GetLength() int {
	return s.length
}

//Returns the capacity of the stack, this
//represents the maximum amount of elements
//that the instance of the stack can hold
func (s *Stack) GetCapacity() int {
	return s.capacity
}

//Returns a string representation of the stack
func (s *Stack) String() string {
	return fmt.Sprintf("{\nData:%v,\nCapacity:%v,\nLength:%v,\n}", s.data, s.capacity, s.length)
}

func (s *Stack) Clone() *Stack {
	s_copy := NewStack(s.capacity)

	for _, element := range s.data {
		s_copy.Push(element)
	}

	return s_copy
}
