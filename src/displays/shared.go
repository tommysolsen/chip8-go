package displays

func memoryLocation(x uint8, y uint8) uint32 {
	return uint32(x)%64 + uint32(y)*64
}
