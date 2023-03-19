package main

import "fmt"

var REG_NARROW = [8]string{"al", "cl", "dl", "bl", "ah", "ch", "dh", "bh"}
var REG_WIDE = [8]string{"ax", "cx", "dx", "bx", "sp", "bp", "si", "di"}
var RM = [8]string{"bx + si", "bx + di", "bp + si", "bp + di", "si", "di", "bp", "bx"}

func GetReg(w byte, reg_lu byte) string {
	if w == 0 {
		return REG_NARROW[reg_lu]
	} else {
		return REG_WIDE[reg_lu]
	}
}

func FmtRm(lu byte, displacement int16) string {
	// format the rm reg content
	rm := RM[lu]

	// if displacement is zero then return reg only
	if displacement == 0 {
		return fmt.Sprintf("[%s]", rm)
	}

	operator := "+"
	if displacement < 0 {
		operator = "-"
	}
	return fmt.Sprintf("[%s %s %d]", rm, operator, Abs(displacement))
}

func WhatAccumulator(m *Memory, w byte) (int16, string) {
	// work out what value is going to/from accumulator
	// work out what accumulator register is
	var dec int16
	var acc string

	if w == 0 {
		dec = m.Get8BitDisplacment()
		acc = "al"
	} else {
		dec = m.Get16BitDisplacement()
		acc = "ax"
	}

	return dec, acc
}
