package stack

import "errors"

type Stack[T any] struct {
	elements []T
}

func (s *Stack[T]) Push(value T) {
	s.elements = append(s.elements, value)
}

func (s *Stack[T]) Pop() (T, error) {
	var nilValue T
	if len(s.elements) == 0 {
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
