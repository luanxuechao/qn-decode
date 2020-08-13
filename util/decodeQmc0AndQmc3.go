package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var maps []int = []int{0x77, 0x48, 0x32, 0x73, 0xDE, 0xF2, 0xC0, 0xC8, 0x95, 0xEC, 0x30, 0xB2, 0x51, 0xC3, 0xE1, 0xA0,
	0x9E, 0xE6, 0x9D, 0xCF, 0xFA, 0x7F, 0x14, 0xD1, 0xCE, 0xB8, 0xDC, 0xC3, 0x4A, 0x67, 0x93, 0xD6,
	0x28, 0xC2, 0x91, 0x70, 0xCA, 0x8D, 0xA2, 0xA4, 0xF0, 8, 0x61, 0x90, 0x7E, 0x6F, 0xA2, 0xE0, 0xEB,
	0xAE, 0x3E, 0xB6, 0x67, 0xC7, 0x92, 0xF4, 0x91, 0xB5, 0xF6, 0x6C, 0x5E, 0x84, 0x40, 0xF7, 0xF3,
	0x1B, 2, 0x7F, 0xD5, 0xAB, 0x41, 0x89, 0x28, 0xF4, 0x25, 0xCC, 0x52, 0x11, 0xAD, 0x43, 0x68, 0xA6,
	0x41, 0x8B, 0x84, 0xB5, 0xFF, 0x2C, 0x92, 0x4A, 0x26, 0xD8, 0x47, 0x6A, 0x7C, 0x95, 0x61, 0xCC,
	0xE6, 0xCB, 0xBB, 0x3F, 0x47, 0x58, 0x89, 0x75, 0xC3, 0x75, 0xA1, 0xD9, 0xAF, 0xCC, 8, 0x73, 0x17,
	0xDC, 0xAA, 0x9A, 0xA2, 0x16, 0x41, 0xD8, 0xA2, 6, 0xC6, 0x8B, 0xFC, 0x66, 0x34, 0x9F, 0xCF, 0x18,
	0x23, 0xA0, 0xA, 0x74, 0xE7, 0x2B, 0x27, 0x70, 0x92, 0xE9, 0xAF, 0x37, 0xE6, 0x8C, 0xA7, 0xBC, 0x62,
	0x65, 0x9C, 0xC2, 8, 0xC9, 0x88, 0xB3, 0xF3, 0x43, 0xAC, 0x74, 0x2C, 0xF, 0xD4, 0xAF, 0xA1, 0xC3, 1,
	0x64, 0x95, 0x4E, 0x48, 0x9F, 0xF4, 0x35, 0x78, 0x95, 0x7A, 0x39, 0xD6, 0x6A, 0xA0, 0x6D, 0x40,
	0xE8, 0x4F, 0xA8, 0xEF, 0x11, 0x1D, 0xF3, 0x1B, 0x3F, 0x3F, 7, 0xDD, 0x6F, 0x5B, 0x19, 0x30, 0x19,
	0xFB, 0xEF, 0xE, 0x37, 0xF0, 0xE, 0xCD, 0x16, 0x49, 0xFE, 0x53, 0x47, 0x13, 0x1A, 0xBD, 0xA4, 0xF1,
	0x40, 0x19, 0x60, 0xE, 0xED, 0x68, 9, 6, 0x5F, 0x4D, 0xCF, 0x3D, 0x1A, 0xFE, 0x20, 0x77, 0xE4, 0xD9,
	0xDA, 0xF9, 0xA4, 0x2B, 0x76, 0x1C, 0x71, 0xDB, 0, 0xBC, 0xFD, 0xC, 0x6C, 0xA5, 0x47, 0xF7, 0xF6, 0,
	0x79, 0x4A, 0x11}

// DecodeQmc0OrQmc3 is  a decode qmc and qmc3 to mp3 adapter
func DecodeQmc0OrQmc3(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf(err.Error())
		return errors.New(err.Error())
	}
	var buffer []byte = data
	for i := 0; i < len(data); i++ {
		buffer[i] = buffer[i] ^ byte(mapL(i))
	}
	strIndex := strings.LastIndex(filename, ".")
	if strIndex == -1 {
		return errors.New("file not expected")
	}
	newFile := filename[0:strIndex] + ".mp3"
	err2 := ioutil.WriteFile(newFile, buffer, 0666)
	if err2 != nil {
		fmt.Printf(err2.Error())
	}
	return nil
}

func mapL(i int) int {
	if i > 0x800 {
		return getI(i % 0x7fff)
	}
	return getI(i)
}
func getI(i int) int {
	return maps[(i*i+80923)%256]
}
