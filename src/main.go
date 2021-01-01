package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args[1:]) == 0 {
        fmt.Println("./emu rom_path")
        return
    }

    state := setupState(os.Args)
    for !state.shouldHalt && state.emulationError == nil {
        emulationCycle(&state)
    }

    if state.emulationError != nil {
        fmt.Print(state.emulationError)
    }
}

func emulationCycle(state *State) {
    instruction := fetchInstruction(state)
    if (state.debug && instruction.opcode != Output) {
        fmt.Println(instruction.opcode, instruction.arguments)
        fmt.Println("Registers: ", state.registers)
        fmt.Println("Stack: ", state.stack)
        fmt.Println()
    }
    executeInstruction(state, instruction)
}
