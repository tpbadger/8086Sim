package main

import "encoding/binary"

func SignExtend8Bits(b byte) []byte {
	lead := (b & 0x80) >> 7
	if lead == 0 {
		return []byte{0x00, b}
	}
	return []byte{0xFF, b}
}

func BytesTodec(decArr []byte) int16 {
	return int16(binary.BigEndian.Uint16(decArr))
}

func Abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}
