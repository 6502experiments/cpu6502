package cpu6502

import "fmt"

type CPU struct {
	// Basic 8=bit registers
	A  byte // Accumulator
	X  byte // Index register (X)
	Y  byte // Index register (Y)
	SP byte // Stack pointer
	P  byte // Status register

	// Program Counter
	PC uint16

	// Cycle count
	Cycles uint64

	// System bus
	bus Bus

	// Instruction table
	instructionTable map[byte]Instruction
}

// ──────────────────────────────────────────────────────────────────────────────
//
//	NewCPU constructor
//
// ──────────────────────────────────────────────────────────────────────────────
func NewCPU(bus Bus) *CPU {
	cpu := &CPU{
		bus: bus,
	}

	// Build instruction table
	cpu.instructionTable = map[byte]Instruction{

		// ──────────────────────────────────────────────────────────────────────
		// NOP (No Operation)
		// ──────────────────────────────────────────────────────────────────────
		// Doesn't actually use an addressing mode, but we can define one
		// (like Immediate) if we want to skip a byte, or NoneAddressing if
		// we just ignore it.
		0x00: {
			Name:   "NOP",
			Cycles: 2,
			Mode: func(c *CPU) uint16 {
				// No real address needed
				return 0
			},
			Execute: func(c *CPU, addr uint16) {
				// Do nothing
			},
		},

		// ──────────────────────────────────────────────────────────────────────
		// LDA Immediate
		// ──────────────────────────────────────────────────────────────────────
		// Loads a value directly from the next byte into A.
		0x01: {
			Name:   "LDA_IMM",
			Cycles: 2,
			Mode:   Immediate,
			Execute: func(c *CPU, addr uint16) {
				// For immediate, 'addr' is the address of the byte in memory
				value := c.bus.Read(addr)
				c.A = value
			},
		},

		// ──────────────────────────────────────────────────────────────────────
		// LDA Zero Page
		// ──────────────────────────────────────────────────────────────────────
		// Reads 8-bit address from next byte, then loads from that zero-page location.
		0x02: {
			Name:   "LDA_ZP",
			Cycles: 3, // Typically more than immediate
			Mode:   ZeroPage,
			Execute: func(c *CPU, addr uint16) {
				c.A = c.bus.Read(addr)
			},
		},
		// ──────────────────────────────────────────────────────────────────────
		// STA Absolute
		// ──────────────────────────────────────────────────────────────────────
		// Stores the accumulator value into an absolute address read from the next 2 bytes.
		0x03: {
			Name:   "STA_ABS",
			Cycles: 4,
			Mode:   Absolute,
			Execute: func(c *CPU, addr uint16) {
				c.bus.Write(addr, c.A)
			},
		},
		// ──────────────────────────────────────────────────────────────────────
		// INC ZeroPage X
		// ──────────────────────────────────────────────────────────────────────
		// Increments whatever value is at (addr + X) in zero page.
		0x04: {
			Name:   "INC_ZPX",
			Cycles: 6,
			Mode:   ZeroPageX,
			Execute: func(c *CPU, addr uint16) {
				val := c.bus.Read(addr)
				val++
				c.bus.Write(addr, val)
			},
		},
	}

	return cpu
}

// ──────────────────────────────────────────────────────────────────────────────
//
//	Step function - executes one instruction
//
// ──────────────────────────────────────────────────────────────────────────────
func (c *CPU) Step() {
	// 1. Fetch opcode from memory at PC
	opcodeAddr := c.PC
	opcode := c.bus.Read(opcodeAddr)

	// 2. Increment PC to point past the opcode
	c.PC++

	// 3. Look up the instruction
	inst, found := c.instructionTable[opcode]
	if !found {
		fmt.Printf("Unknown opcode 0x%02X at 0x%04X. Treating as NOP.\n", opcode, opcodeAddr)
		c.Cycles += 2
		return
	}

	// 4. Determine the effective address using the addressing mode function
	//    (some modes might not need to do anything, but we call it anyway for consistency).
	effectiveAddr := inst.Mode(c)

	// 5. Execute the instruction, passing the CPU and the effective address
	inst.Execute(c, effectiveAddr)

	// 6. Update cycle count
	c.Cycles += uint64(inst.Cycles)
}
