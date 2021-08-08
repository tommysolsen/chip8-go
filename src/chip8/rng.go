package chip8

import "math/rand"

type rngGenerator struct {
}

func (r rngGenerator) GetRandom() uint8 {
	return uint8(rand.Intn(8))
}

type rngGeneratorMock struct {
	Value uint8
}

func (r rngGeneratorMock) GetRandom() uint8 {
	return r.Value
}
