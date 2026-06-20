package stack

import (
	"errors"
	"fmt"
)

type Stack[T any] struct {
	elements []T
}

func (s *Stack[T]) Push(value T) {
	s.elements = append(s.elements, value)
}

func (s *Stack[T]) Pop() (T, error) {
	var nilValue T
	if s.IsEmpty() {
		return nilValue, errors.New("stack is empty")
	}

	i := len(s.elements) - 1
	element := s.elements[i]

	s.elements[i] = nilValue
	s.elements = s.elements[:i]

	return element, nil
}

func (s *Stack[T]) Peek() (T, error) {
	var nilValue T
	if len(s.elements) == 0 {
		return nilValue, errors.New("stack is empty")
	}

	i := len(s.elements) - 1
	return s.elements[i], nil
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.elements)
}

func (s *Stack[T]) Get(index int) (T, error) {
	if index < 0 || index >= s.Size() {
		var nilValue T
		return nilValue, fmt.Errorf("index %d out of bounds for stack size %d", index, len(s.elements))
	}

	return s.elements[index], nil
}
