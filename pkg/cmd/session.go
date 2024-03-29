package cmd

import (
	"fmt"
	"os"
	"strconv"
)

//				Structural things				//

// S interface describes what session can do
type S interface {
	Input(expression []string)
	PrintStack()
	PrintSymbolTable()
	Clear()

	ConvertOutputAndVariables(variable1, variable2 string) (int, int)
	AssignValueToAVariable(variableName, variableValue string)
}

type Session struct {
	Stack     *Stack
	Variables *SymbolTable
}

func NewSession() *Session {
	s := Session{}
	s.Stack = NewStack()
	s.Variables = NewSymbolTable()

	return &s
}

//	Method implementation					//
//
// Input iterates over the expression tokens to do the interpretation
/*
	INPUT(expression):
	1. for token := expression
	2. 		if token is number
	3.			stack.push(token)
	4.		if token is operator
	5.			operand1 = stack.pop()
	6.			operand2 = stack.pop()
	7.			switch (token)
	8.				case "+": stack.push(operand1 + operand2)
	9.				case "-": stack.push(operand1 - operand2)
	10.				case "*": stack.push(operand1 * operand2)
	11.				case "/": stack.push(operand1 / operand2)
	12.				case "=":
	13.					symbol_table.insert(operand2, operand1)
	14.		if token.size == 1 and token in 'A' to 'Z'
	15.			stack.push(token)
*/
func (s *Session) Input(expression []string) {
	for _, token := range expression {
		// If conversion error is nil it is a value, and we can put it in the Stack.
		// Else if it is not nil, it is an operator, and we have to do some arithmetics.
		_, err := strconv.Atoi(token)
		if err == nil {
			s.Stack.Push(token)
			s.PrintStack()
		} else {
			if IsValidOperator(token) {
				stringValue1, exists := s.Stack.Pop()
				if !exists {
					Red.Println("wrong expression, no operands for the following operator:", token)
					return
				}

				stringValue2, exists := s.Stack.Pop()
				if !exists {
					Red.Println("wrong expression, no operands for the following operator:", token)
					return
				}

				switch token {
				case "+":
					value1, value2 := s.ConvertOutputAndVariables(stringValue1, stringValue2)

					result := value1 + value2
					s.Stack.Push(strconv.Itoa(result))
					s.PrintStack()
				case "-":
					value1, value2 := s.ConvertOutputAndVariables(stringValue1, stringValue2)

					result := value1 - value2
					s.Stack.Push(strconv.Itoa(result))
					s.PrintStack()
				case "*":
					value1, value2 := s.ConvertOutputAndVariables(stringValue1, stringValue2)

					result := value1 * value2
					s.Stack.Push(strconv.Itoa(result))
					s.PrintStack()
				case "/":
					value1, value2 := s.ConvertOutputAndVariables(stringValue1, stringValue2)

					// Well, we of course can make a Stack of float/double values,
					// but I believe taking the whole part of division will be
					// illustrative enough.
					result := value1 / value2
					s.Stack.Push(strconv.Itoa(result))
					s.PrintStack()
				case "=":
					// do an assignment here
					s.AssignValueToAVariable(stringValue2, stringValue1)
					s.PrintStack()
				default:
					Red.Println("the following operator is not supported:", token)
					return
				}
			} else if len(token) == 1 && token[0] >= 'A' && token[0] <= 'Z' {
				s.Stack.Push(token)
			} else {
				Red.Println("this kind of input is not supported", token)
				Red.Println("a variable can only be in the range of A-Z")
				return
			}
		}
	}
}

// PrintStack prints the current content of the stack
func (s *Session) PrintStack() {
	// Result is calculated by unpacking the Stack and doing all the operations
	tempStack := *s.Stack
	fmt.Printf("[ ")

	for !tempStack.IsEmpty() {
		pop, exists := tempStack.Pop()
		if !exists {
			break
		}

		fmt.Printf(pop + " ")

	}
	fmt.Printf("]\n")
}

// PrintSymbolTable prints the current variables in the state
func (s *Session) PrintSymbolTable() {
	fmt.Print(s.Variables.String())
}

// Clear purges the session, reinstantiating the stack and symbol table
func (s *Session) Clear() {
	s.Stack = NewStack()
	s.Variables = NewSymbolTable()
	err := os.Remove("session.json")
	if err != nil {
		panic("cannot rm session.json")
	}
}

//	Utility methods				//
//
// ConvertOutputAndVariables validates two values popped from the stack and
// identifies if it is a variable or a number
func (s *Session) ConvertOutputAndVariables(variable1, variable2 string) (int, int) {
	value1, convErr := strconv.Atoi(variable1)
	if convErr != nil {
		variableStringValue, variableExists := s.Variables.Get(variable1)
		if !variableExists {
			Red.Println("there is no such variable:", variable1)
			return 0, 0
		}

		variableIntValue, convErr := strconv.Atoi(variableStringValue)
		if convErr != nil {
			Red.Println("problem with variable value conversion from string to int")
			return 0, 0
		}

		value1 = variableIntValue
	}
	value2, convErr := strconv.Atoi(variable2)
	if convErr != nil {
		variableStringValue, variableExists := s.Variables.Get(variable2)
		if !variableExists {
			Red.Println("there is no such variable:", variable2)
			os.Exit(1)
		}

		variableIntValue, convErr := strconv.Atoi(variableStringValue)
		if convErr != nil {
			Red.Println("problem with variable value conversion from string to int")
			return 0, 0
		}

		value2 = variableIntValue
	}

	return value1, value2
}

// AssignValueToAVariable validates the variable name and value and instantiates it in the symbol table
func (s *Session) AssignValueToAVariable(variableName, variableValue string) {
	_, convErr := strconv.Atoi(variableName)
	if convErr == nil {
		// Successful conversion cannot appear here, it means we did something wrong,
		// as variableName should not be integer. So we throw an error message.
		Red.Println("wrong assignment order")
		return
	}

	_, convErr = strconv.Atoi(variableValue)
	if convErr != nil {
		Red.Println("wrong assignment argument:", variableValue)
		return
	}

	s.Variables.Insert(variableName, variableValue)
}
