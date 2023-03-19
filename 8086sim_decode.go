package main

import (
	"fmt"
)

func DecodeInstructions(m *Memory) {
	// run through memory and decode instruction stream
	// print output to std out

	for m.loc < len(m.a) {

		currentByte := m.GetByte()

		if currentByte>>2 == 0x22 {
			// MOV
			GenericEitherToEither(m, "mov")
		} else if currentByte>>4 == 0x0B {
			MovImmediateToRegister(m)
		} else if currentByte>>1 == 0x50 {
			MovMemoryToAccumulator(m, "mov")
		} else if currentByte>>1 == 0x51 {
			MovMemoryFromAccumulator(m, "mov")
		} else if currentByte>>1 == 0x63 {
			MovImmediateToEither(m)
			// ARITHMETIC
		} else if currentByte>>2 == 0x20 {
			ArithmeticImmediateToEither(m)
			// ADD
		} else if currentByte>>2 == 0x00 {
			GenericEitherToEither(m, "add")
		} else if currentByte>>1 == 0x02 {
			ArithmeticDataToAccumulator(m, "add")
			// SUB
		} else if currentByte>>2 == 0x0A {
			GenericEitherToEither(m, "sub")
		} else if currentByte>>1 == 0x16 {
			ArithmeticDataToAccumulator(m, "sub")
			// CMP
		} else if currentByte>>2 == 0x0E {
			GenericEitherToEither(m, "cmp")
		} else if currentByte>>1 == 0x1E {
			ArithmeticDataToAccumulator(m, "cmp")
		} else {
			panic(fmt.Errorf("cant decode byte %08b", currentByte))
		}

	}
}

func GenericEitherToEither(m *Memory, op string) {
	currentByte := m.GetByteIncr()

	var reg string
	var rm string

	d := currentByte & 0x02
	w := currentByte & 0x01

	currentByte = m.GetByteIncr()

	mod := (currentByte & 0xC0) >> 6

	rm_lu := currentByte & 0x07
	reg_lu := (currentByte & 0x38) >> 3
	reg = GetReg(w, reg_lu)

	switch mod {
	case 0:
		// no displacement
		// unless rm_lu == 6 (110) then 16 bit
		if rm_lu == 6 {
			rm = fmt.Sprintf("[%d]", m.Get16BitDisplacement())
		} else {
			rm = fmt.Sprintf("[%s]", RM[rm_lu])
		}

	case 1:
		// 8 bit displacement
		dec := m.Get8BitDisplacment()
		rm = FmtRm(rm_lu, dec)

	case 2:
		// 16 bit displacement
		dec := m.Get16BitDisplacement()
		rm = FmtRm(rm_lu, dec)

	case 3:
		// no displacement
		// reg = rm
		rm = GetReg(w, rm_lu)
	}

	if d == 0 {
		fmt.Printf("%s %s, %s\n", op, rm, reg)
	} else {
		fmt.Printf("%s %s, %s\n", op, reg, rm)
	}
}

func GenericFromAccumulator(m *Memory, op string) {
	currentByte := m.GetByteIncr()

	w := currentByte & 0x01
	dec, acc := WhatAccumulator(m, w)
	fmt.Printf("%s [%d], %s\n", op, dec, acc)
}

func MovMemoryToAccumulator(m *Memory, op string) {
	currentByte := m.GetByteIncr()

	w := currentByte & 0x01
	dec, acc := WhatAccumulator(m, w)
	fmt.Printf("%s %s, [%d]\n", op, acc, dec)
}

func MovMemoryFromAccumulator(m *Memory, op string) {
	currentByte := m.GetByteIncr()

	w := currentByte & 0x01
	dec, acc := WhatAccumulator(m, w)
	fmt.Printf("%s [%d], %s\n", op, dec, acc)
}

func MovImmediateToRegister(m *Memory) {
	currentByte := m.GetByteIncr()

	var reg string
	var dec int16

	reg_lu := currentByte & 0x07
	w := (currentByte & 0x08) >> 3

	if w == 0 {
		reg = REG_NARROW[reg_lu]
		dec = m.Get8BitDisplacment()
	} else {
		reg = REG_WIDE[reg_lu]
		dec = m.Get16BitDisplacement()
	}

	fmt.Printf("mov %s, %d\n", reg, dec)
}

func MovImmediateToEither(m *Memory) {
	// i know this is basically carbon copy of GenericEitherToEither
	// but this has its own nuances and as its a seperate instruction
	// am happy with having some repeated code for the sake of clarity
	currentByte := m.GetByteIncr()

	var rm string
	var dec int16
	var decs string

	w := currentByte & 0x01

	currentByte = m.GetByteIncr()

	mod := currentByte & 0xC0 >> 6
	rm_lu := currentByte & 0x07

	switch mod {
	case 0:
		// no displacement
		// unless rm_lu == 6 (110) then 16 bit
		if rm_lu == 6 {
			rm = fmt.Sprintf("[%d]", m.Get16BitDisplacement())
		} else {
			rm = fmt.Sprintf("[%s]", RM[rm_lu])
		}

	case 1:
		// 8 bit displacement
		dec := m.Get8BitDisplacment()
		rm = FmtRm(rm_lu, dec)

	case 2:
		// 16 bit displacement
		dec := m.Get16BitDisplacement()
		rm = FmtRm(rm_lu, dec)

	case 3:
		// no displacement
		// reg = rm
		rm = GetReg(w, rm_lu)
	}

	if w == 0 {
		// 8 bit data
		dec = m.Get8BitDisplacment()
		decs = fmt.Sprintf("word %d", dec)
	} else {
		// 16 bit data
		dec = m.Get16BitDisplacement()
		decs = fmt.Sprintf("byte %d", dec)
	}

	fmt.Printf("mov %s, %s\n", rm, decs)
}

func ArithmeticDataToAccumulator(m *Memory, op string) {
	currentByte := m.GetByteIncr()

	w := currentByte & 0x01
	dec, acc := WhatAccumulator(m, w)
	fmt.Printf("%s %s, %d\n", op, acc, dec)
}

func ArithmeticDataFromAccumulator(m *Memory, op string) {
	currentByte := m.GetByteIncr()

	w := currentByte & 0x01
	dec, acc := WhatAccumulator(m, w)
	fmt.Printf("%s %s, %d\n", op, acc, dec)
}

func ArithmeticImmediateToEither(m *Memory) {
	currentByte := m.GetByteIncr()

	var rm string
	var dec int16
	var decs string
	var ops string

	sw := currentByte & 0x02
	w := currentByte & 0x01

	currentByte = m.GetByteIncr()

	mod := currentByte & 0xC0 >> 6
	rm_lu := currentByte & 0x07
	op := (currentByte & 0x38) >> 3

	switch mod {
	case 0:
		// no displacement
		// unless rm_lu == 6 (110) then 16 bit
		if rm_lu == 6 {
			rm = fmt.Sprintf("[%d]", m.Get16BitDisplacement())
		} else {
			rm = fmt.Sprintf("[%s]", RM[rm_lu])
		}

	case 1:
		// 8 bit displacement
		dec := m.Get8BitDisplacment()
		rm = FmtRm(rm_lu, dec)

	case 2:
		// 16 bit displacement
		dec := m.Get16BitDisplacement()
		rm = FmtRm(rm_lu, dec)

	case 3:
		// no displacement
		// reg = rm
		rm = GetReg(w, rm_lu)
	}

	switch op {
	case 0:
		// 000
		ops = "add"
	case 5:
		// 101
		ops = "sub"
	case 7:
		// 111
		ops = "cmp"
	}

	if sw == 1 {
		dec = m.Get16BitDisplacement()
		decs = fmt.Sprintf("byte %d", dec)
	} else {
		dec = m.Get8BitDisplacment()
		decs = fmt.Sprintf("word %d", dec)
	}

	fmt.Printf("%s %s, %s\n", ops, rm, decs)
}
