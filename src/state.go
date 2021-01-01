package main

import (
    "errors"
    "io/ioutil"
    "os"
)

type State struct {
    debug               bool
    emulationError      error
    memory              []byte
    programCounter      uint16
    registers           [8]uint16
    shouldHalt          bool
    stack               Stack
}

func setupState(args []string) State {
    state := State {
        debug: false,
        emulationError: nil,
        programCounter: 0,
        shouldHalt: false,
        stack: Stack {
            stack: []uint16{},
        },
    }

    romPath := ""
    for _, arg := range os.Args[1:] {
        if arg == "-d" {
            state.debug = true
        } else if _, err := os.Stat(arg); !os.IsNotExist(err) {
            romPath = arg
        }
    }

    if romPath == "" {
        state.emulationError = errors.New("Error: Rom does not exists.")
        state.shouldHalt = true
        return state
    }

    rom, err := ioutil.ReadFile(romPath)
    if err != nil {
        state.emulationError = errors.New("Error: Unable to load rom.")
        state.shouldHalt = true
        return state
    }

    freeMemory := make([]byte, 0xFFFE - len(rom))
    state.memory = append(rom, freeMemory...)

    return state
}
