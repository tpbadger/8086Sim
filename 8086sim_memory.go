package main

import (
	"fmt"
	"io"
	"os"
)

type Memory struct {
	a   []byte
	loc int
}

func (m Memory) PrintAllMemory() {
	// for debugging
	fmt.Printf("%08b\n", m.a)
}

func (m Memory) PrintCurrentByte() {
	// for debugging
	fmt.Printf("%08b\n", m.a[m.loc])
}

func (m *Memory) GetByteIncr() byte {
	// returns the byte at current memory array loc
	// increments loc ready for next call
	b := m.a[m.loc]
	m.loc += 1
	return b
}

func (m *Memory) GetByte() byte {
	// returns the byte at current memory array loc
	return m.a[m.loc]
}

func (m *Memory) Get16BitDisplacement() int16 {
	lsb := m.GetByteIncr()
	msb := m.GetByteIncr()
	decArr := []byte{msb, lsb}
	return BytesTodec(decArr)
}

func (m *Memory) Get8BitDisplacment() int16 {
	// sign extend 8 bits
	lsb := m.GetByteIncr()
	decArr := SignExtend8Bits(lsb)
	return BytesTodec(decArr)
}

func LoadFileIntoMemory(fpth string, m *Memory) error {

	file, err := os.Open(fpth)

	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	if len(bytes) > 1000 {
		return fmt.Errorf("instruction stream of %d bytes exceeds 8086 memory of 1mb", len(bytes))
	}

	m.a = bytes

	return nil
}
