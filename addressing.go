package cpu6502

// ──────────────────────────────────────────────────────────────────────────────
//  Addressing Mode
// ──────────────────────────────────────────────────────────────────────────────
// Each addressing mode returns an effective address (the location in memory
// to be read from or written to). Some addressing modes might not need to
// read memory (e.g., immediate mode), but we'll standardize by returning
// an address. The CPU or the instruction can then decide how to handle it.

type AddressingMode func(c *CPU) uint16

// Immediate returns the address of the next byte in the program memory (PC).
func Immediate(c *CPU) uint16 {
	addr := c.PC
	c.PC++
	return addr
}

// ZeroPage uses the next byte as an 8-bit address in page zero (0x00xx).
// The effective address is thus 0x00XX.
func ZeroPage(c *CPU) uint16 {
	// Fetch a single byte
	offset := c.bus.Read(c.PC)
	c.PC++
	return uint16(offset)
}

// Absolute uses the next two bytes as a 16-bit address.
func Absolute(c *CPU) uint16 {
	low := c.bus.Read(c.PC)
	high := c.bus.Read(c.PC + 1)
	addr := uint16(low) | (uint16(high) << 8)
	c.PC += 2
	return addr
}

// ZeroPageX is a zero page address offset by the X register.
func ZeroPageX(c *CPU) uint16 {
	base := c.bus.Read(c.PC)
	c.PC++
	return uint16(base+c.X) & 0x00FF
}

// AbsoluteX is a 16-bit address plus the X register.
func AbsoluteX(c *CPU) uint16 {
	low := c.bus.Read(c.PC)
	high := c.bus.Read(c.PC + 1)
	addr := uint16(low) | (uint16(high) << 8)
	c.PC += 2
	return addr + uint16(c.X)
}
