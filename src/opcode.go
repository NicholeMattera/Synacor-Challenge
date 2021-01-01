package main

type Opcode int
const (
    Halt Opcode = iota
    Set
    Push
    Pop
    EqualTo
    GreaterThan
    JumpTo
    JumpToIfNotZero
    JumpToIfZero
    Add
    Multiply
    Modulous
    And
    Or
    Not
    ReadMemory
    WriteMemory
    Call
    Return
    Output
    Input
    NoOperation
)

func (o Opcode) String() string {
    return [...]string {
        "HALT",
        "SET",
        "PUSH",
        "POP",
        "EQ",
        "GT",
        "JMP",
        "JT",
        "JF",
        "ADD",
        "MULT",
        "MOD",
        "AND",
        "OR",
        "NOT",
        "RMEM",
        "WMEM",
        "CALL",
        "RET",
        "OUT",
        "IN",
        "NOOP",
    }[o]
}
