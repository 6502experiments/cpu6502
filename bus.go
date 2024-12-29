package cpu6502

// Bus interface
type Bus interface {
	Read(addr uint16) byte
	Write(addr uint16, data byte)
}
