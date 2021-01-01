package main

import (
    "errors"
    "fmt"
    "github.com/eiannone/keyboard"
)

type Instruction struct {
    opcode      Opcode
    arguments   []Number
}

func fetchInstruction(state *State) Instruction {
    if state.programCounter > uint16(32767) {
        return Instruction {
            opcode: Halt,
            arguments: []Number{},
        }
    }

    instruction := Instruction {
        opcode: Opcode(state.memory[state.programCounter]),
        arguments: []Number{},
    }
    programCounterIncrement := uint16(2)

    switch instruction.opcode {
        case EqualTo:
            fallthrough
        case GreaterThan:
            fallthrough
        case Add:
            fallthrough
        case Multiply:
            fallthrough
        case Modulous:
            fallthrough
        case And:
            fallthrough
        case Or: 
            instruction.arguments = append([]Number {
                parseNumber(state.memory[state.programCounter + 6:state.programCounter + 8], state),
            }, instruction.arguments...)
            programCounterIncrement += uint16(2)
            fallthrough 

        case Set:
            fallthrough 
        case JumpToIfNotZero:
            fallthrough 
        case JumpToIfZero:
            fallthrough 
        case Not:
            fallthrough 
        case ReadMemory:
            fallthrough 
        case WriteMemory:
            instruction.arguments = append([]Number {
                parseNumber(state.memory[state.programCounter + 4:state.programCounter + 6], state),
            }, instruction.arguments...)
            programCounterIncrement += uint16(2)
            fallthrough 

        case Push:
            fallthrough 
        case Pop:
            fallthrough 
        case JumpTo:
            fallthrough 
        case Call:
            fallthrough 
        case Output:
            fallthrough 
        case Input:
            instruction.arguments = append([]Number {
                parseNumber(state.memory[state.programCounter + 2:state.programCounter + 4], state),
            }, instruction.arguments...)
            programCounterIncrement += uint16(2)
            break 
    }

    state.programCounter += programCounterIncrement

    return instruction
}

func executeInstruction(state *State, instruction Instruction) {
    switch instruction.opcode {
        // HALT
        // Stop execution and terminate the program.
        case Halt:
            state.shouldHalt = true
            break

        // SET a b
        // set register <a> to the value of <b>.
        case Set:
            a := instruction.arguments[0]
            b := instruction.arguments[1]

            if a.isRegister {
                state.registers[a.registerIndex] = b.value
            }
            break

        // PUSH a
        // Push <a> onto the stack.
        case Push:
            a := instruction.arguments[0]

            state.stack.push(a.value)
            break

        // POP a
        // Remove the top element from the stack and write it into <a>; empty stack = error.
        case Pop:
            a := instruction.arguments[0]

            if  a.isRegister && state.stack.length() != 0 {
                state.registers[a.registerIndex] = state.stack.pop().value
            }
            break

        // EQ
        // Set <a> to 1 if <b> is equal to <c>; set it to 0 otherwise.
        case EqualTo:
            a := instruction.arguments[0]
            b := instruction.arguments[1]
            c := instruction.arguments[2]

            if a.isRegister {
                if b.value == c.value {
                    state.registers[a.registerIndex] = 1
                } else {
                    state.registers[a.registerIndex] = 0
                }
            }
            break

        // GT a b c
        // Set <a> to 1 if <b> is greater than <c>; set it to 0 otherwise.
        case GreaterThan:
            a := instruction.arguments[0]
            b := instruction.arguments[1]
            c := instruction.arguments[2]

            if a.isRegister {
                if b.value > c.value {
                    state.registers[a.registerIndex] = 1
                } else {
                    state.registers[a.registerIndex] = 0
                }
            }
            break

        // JMP a
        // Jump to <a>.
        case JumpTo:
            a := instruction.arguments[0]

            state.programCounter = a.toAddress()
            break

        // JT a b
        // If <a> is nonzero, jump to <b>.
        case JumpToIfNotZero:
            a := instruction.arguments[0]
            b := instruction.arguments[1]

            if a.value != 0 {
                state.programCounter = b.toAddress()
            } 
            break

        // JF a b
        // If <a> is zero, jump to <b>.
        case JumpToIfZero:
            a := instruction.arguments[0]
            b := instruction.arguments[1]
            
            if a.value == 0 {
                state.programCounter = b.toAddress()
            } 
            break

        // ADD a b c
        // Assign into <a> the sum of <b> and <c>. (Modulo 32768)
        case Add:
            a := instruction.arguments[0]
            b := instruction.arguments[1]
            c := instruction.arguments[2]

            if a.isRegister {
                state.registers[a.registerIndex] = (b.value + c.value) % 32768
            }
            break

        // MULT a b c
        // Store into <a> the product of <b> and <c>. (Modulo 32768)
        case Multiply:
            a := instruction.arguments[0]
            b := instruction.arguments[1]
            c := instruction.arguments[2]

            if a.isRegister {
                state.registers[a.registerIndex] = (b.value * c.value) % 32768
            }
            break

        // MOD a b c
        // Store into <a> the remainder of <b> divided by <c>.
        case Modulous:
            a := instruction.arguments[0]
            b := instruction.arguments[1]
            c := instruction.arguments[2]

            if a.isRegister {
                state.registers[a.registerIndex] = b.value % c.value
            }
            break

        // AND a b c
        // Stores into <a> the bitwise and of <b> and <c>.
        case And:
            a := instruction.arguments[0]
            b := instruction.arguments[1]
            c := instruction.arguments[2]

            if a.isRegister {
                state.registers[a.registerIndex] = b.value & c.value
            }
            break

        // OR a b c
        // Stores into <a> the bitwise or of <b> and <c>.
        case Or:
            a := instruction.arguments[0]
            b := instruction.arguments[1]
            c := instruction.arguments[2]

            if a.isRegister {
                state.registers[a.registerIndex] = b.value | c.value
            }
            break

        // NOT a b
        // Stores 15-bit bitwise inverse of <b> in <a>.
        case Not:
            a := instruction.arguments[0]
            b := instruction.arguments[1]

            if a.isRegister {
                state.registers[a.registerIndex] = 0x7FFF ^ b.value
            }
            break

        // RMEM a b
        // Read memory at address <b> and write it to <a>.
        case ReadMemory:
            a := instruction.arguments[0]
            b := instruction.arguments[1]

            if a.isRegister {
                address := b.toAddress()
                state.registers[a.registerIndex] = parseNumber(state.memory[address:address + 2], state).value
            }
            break

        // WMEM a b
        // Write the value from <b> into memory at address <a>.
        case WriteMemory:
            a := instruction.arguments[0]
            b := instruction.arguments[1]

            address := a.toAddress()
            bytes := b.toBytes()
            state.memory[address] = bytes[0]
            state.memory[address + 1] = bytes[1]
            break

        // CALL a
        // Write the address of the next instruction to the stack and jump to <a>.
        case Call:
            a := instruction.arguments[0]
            
            state.stack.push(state.programCounter / 2)
            state.programCounter = a.toAddress()
            break

        // RET
        // Remove the top element from the stack and jump to it; empty stack = halt.
        case Return:
            if state.stack.length() == 0 {
                state.shouldHalt = true
            } else {
                state.programCounter = state.stack.pop().toAddress()
            }
            break
        
        // OUT a
        // Write the character represented by ascii code <a> to the terminal.
        case Output:
            fmt.Print(string(instruction.arguments[0].value))
            break

        // IN a
        // Read a character from the terminal and write its ascii code to <a>; it can be assumed that once input starts, it will continue until a newline is encountered; this means that you can safely read whole lines from the keyboard and trust that they will be fully read.
        case Input:
            a := instruction.arguments[0]

            char, key, _ := keyboard.GetSingleKey()
            fmt.Print(string(char))

            if key == keyboard.KeyEnter {
                state.registers[a.registerIndex] = 0x0A
            } else if key == keyboard.KeySpace {
                state.registers[a.registerIndex] = 0x20
            } else {
                state.registers[a.registerIndex] = uint16(char)
            }
            break

        // NOOP
        // No operation.
        case NoOperation:
            break

        default:
            state.emulationError = errors.New("Error: Invalid opcode!")
            break
    }
}
