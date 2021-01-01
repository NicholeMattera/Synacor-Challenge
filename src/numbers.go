package main

import (
    "encoding/binary"
    "errors"
    "fmt"
)

type Number struct {
    isRegister      bool
    value           uint16
    registerIndex   uint16
}

func (n Number) String() string {
    if n.isRegister {
        return fmt.Sprintf("Register #%d - %d", n.registerIndex, n.value)
    } else {
        return fmt.Sprintf("%d", n.value)
    }
}

func parseNumber(memory []byte, state *State) Number {
    value := binary.LittleEndian.Uint16(memory)

    if value <= 32767 {
        return Number {
            isRegister: false,
            value: value,
            registerIndex: 0,
        }
    } else if value >= 32768 && value <= 32775 {
        index := value - 32768
        return Number {
            isRegister: true,
            value: state.registers[index],
            registerIndex: index,
        }
    } else {
        state.emulationError = errors.New("Error: Invalid Number")
        return Number {
            isRegister: false,
            value: 0,
            registerIndex: 0,
        }
    }
}

func (n Number) toAddress() uint16 {
    return n.value * 2
}

func (n Number) toBytes() []byte {
    bytes := make([]byte, 2)
    binary.LittleEndian.PutUint16(bytes, n.value)
    return bytes
}
