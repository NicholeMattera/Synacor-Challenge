package main

import (
    "fmt"
)

type Stack struct {
    stack   []uint16
}

func (s Stack) String() string {
    return fmt.Sprintf("%v", s.stack)
}

func (s *Stack) push(data uint16) {
    s.stack = append([]uint16{ data }, s.stack...)
}

func (s *Stack) pop() Number {
    data := Number {
        isRegister: false,
        value: s.stack[0],
        registerIndex: 0,
    }
    s.stack = s.stack[1:]
    return data
}

func (s Stack) length() uint16 {
    return uint16(len(s.stack))
}
