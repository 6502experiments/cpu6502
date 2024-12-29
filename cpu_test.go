package cpu6502_test

import (
	"testing"

	cpu "github.com/6502experiments/cpu6502"
)

// Define a bus to connect the CPU to the memory
type Memory struct {
	memory [0x10000]byte
}

func (m *Memory) Read(addr uint16) byte {
	return m.memory[addr]
}

func (m *Memory) Write(addr uint16, value byte) {
	m.memory[addr] = value
}

func NewMemory() *Memory {
	return &Memory{}
}

func TestLDAImmediate(t *testing.T) {
	mem := NewMemory()

	// Program:
	//   0x0000: 0x01 (LDA_IMM), 0x42 (immediate value)
	mem.Write(0x0000, 0x01) // LDA_IMM opcode
	mem.Write(0x0001, 0x42) // immediate operand

	myCPU := cpu.NewCPU(mem)

	// Execute one instruction
	myCPU.Step()

	// We expect the accumulator to have 0x42
	if myCPU.A != 0x42 {
		t.Errorf("LDA_IMM failed. Got A = 0x%02X; want 0x42", myCPU.A)
	}

	// We expect the cycle count to increase by 2 (based on the instruction definition)
	if myCPU.Cycles != 2 {
		t.Errorf("Cycles mismatch after LDA_IMM. Got %d; want 2", myCPU.Cycles)
	}
}

func TestLDAZeroPage(t *testing.T) {
	mem := NewMemory()

	// Program:
	//   0x0000: 0x02 (LDA_ZP), 0x10 (the zero-page address)
	//   0x0010: 0x99 (the value we expect to load)
	mem.Write(0x0000, 0x02) // LDA_ZP opcode
	mem.Write(0x0001, 0x10) // zero-page address
	mem.Write(0x0010, 0x99) // the value we expect to load

	myCPU := cpu.NewCPU(mem)

	// Execute one instruction
	myCPU.Step()

	// We expect the accumulator to have 0x99
	if myCPU.A != 0x99 {
		t.Errorf("LDA_ZP failed. Got A = 0x%02X; want 0x99", myCPU.A)
	}

	// LDA_ZP cycles are set to 3 in the example
	if myCPU.Cycles != 3 {
		t.Errorf("Cycles mismatch after LDA_ZP. Got %d; want 3", myCPU.Cycles)
	}
}

func TestSTAAbsolute(t *testing.T) {
	mem := NewMemory()

	// Program:
	//   0x0000: 0x03 (STA_ABS), then address 0x200 (two bytes: 0x00, 0x02)
	mem.Write(0x0000, 0x03) // STA_ABS opcode
	mem.Write(0x0001, 0x00) // low byte of address
	mem.Write(0x0002, 0x02) // high byte of address

	myCPU := cpu.NewCPU(mem)
	// Put something in accumulator to store
	myCPU.A = 0x55

	// Execute one instruction
	myCPU.Step()

	// The absolute address is 0x0200
	val := mem.Read(0x0200)
	if val != 0x55 {
		t.Errorf("STA_ABS failed. Memory at 0x200 = 0x%02X; want 0x55", val)
	}

	// STA_ABS cycles are set to 4 in the example
	if myCPU.Cycles != 4 {
		t.Errorf("Cycles mismatch after STA_ABS. Got %d; want 4", myCPU.Cycles)
	}
}

func TestNOP(t *testing.T) {
	mem := NewMemory()

	// Program:
	//   0x0000: 0x00 (NOP)
	mem.Write(0x0000, 0x00) // NOP opcode

	myCPU := cpu.NewCPU(mem)

	// Execute one instruction
	myCPU.Step()

	// After NOP, we expect no change to registers,
	// but cycles should have increased by 2 (based on the definition).
	if myCPU.Cycles != 2 {
		t.Errorf("Cycles mismatch after NOP. Got %d; want 2", myCPU.Cycles)
	}

	// We can also check PC advanced by 1.
	if myCPU.PC != 0x0001 {
		t.Errorf("Program Counter mismatch after NOP. Got PC = 0x%04X; want 0x0001", myCPU.PC)
	}
}
