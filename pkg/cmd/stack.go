package cmd

// The most basic implementation of stack data structure //

func NewStack() *Stack {
	return &Stack{Values: make([]string, 0)}
}

type Stack struct {
	Values []string `json:"values"`
}

// IsEmpty checks if there are no elements in the stack
func (s *Stack) IsEmpty() bool {
	return len(s.Values) == 0
}

// Push pushes the value in the stack
/*
	PUSH(value):
	1. stack.push_back(value)
*/

func (s *Stack) Push(value string) {
	s.Values = append(s.Values, value) // Simply append the new value to the end of the stack
}

// Pop pops out a value from the stack by taking a subslice
/*
	POP():
	1. if stack.size == 0
	2. 		return "", false
	3. else
	4. 		i = stack.size - 1
	5.		el = stack[i]
	6.		stack = stack[0:i]
	7.		return el, true
*/
func (s *Stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		// Take the last element out and crop the stack by one element [:index]
		index := len(s.Values) - 1
		element := (s.Values)[index]
		s.Values = (s.Values)[:index]

		return element, true
	}
}
