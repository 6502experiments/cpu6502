package cpu6502

// ──────────────────────────────────────────────────────────────────────────────
//  Instruction struct
// ──────────────────────────────────────────────────────────────────────────────
// Each instruction knows:
//  - Name: For debugging/logging
//  - Cycles: Base cycles
//  - Mode: Addressing mode (function) used to find the effective address
//  - Execute: The actual logic that uses CPU state and the effective address
type Instruction struct {
	Name    string
	Cycles  int
	Mode    AddressingMode
	Execute func(c *CPU, addr uint16)
}
