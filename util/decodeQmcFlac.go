package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var seedMap [8][7]int = [8][7]int{
	{0x4a, 0xd6, 0xca, 0x90, 0x67, 0xf7, 0x52},
	{0x5e, 0x95, 0x23, 0x9f, 0x13, 0x11, 0x7e},
	{0x47, 0x74, 0x3d, 0x90, 0xaa, 0x3f, 0x51},
	{0xc6, 0x09, 0xd5, 0x9f, 0xfa, 0x66, 0xf9},
	{0xf3, 0xd6, 0xa1, 0x90, 0xa0, 0xf7, 0xf0},
	{0x1d, 0x95, 0xde, 0x9f, 0x84, 0x11, 0xf4},
	{0x0e, 0x74, 0xbb, 0x90, 0xbc, 0x3f, 0x92},
	{0x00, 0x09, 0x5b, 0x9f, 0x62, 0x66, 0xa1}}

// QmcFlac2MP3Inf model
type QmcFlac2MP3Inf struct {
	Ret   int
	X     int
	Y     int
	Dx    int
	Index int
}

// DecodeQmcFlac  qmcflac to mp3
func DecodeQmcFlac(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf(err.Error())
		return errors.New(err.Error())
	}
	var buffer []byte = data
	x := -1
	y := 8
	dx := 1
	index := -1
	for i := 0; i < len(data); i++ {
		result := qmcFlac2MP3(x, y, dx, index)
		buffer[i] = byte(result.Ret) ^ buffer[i]
		x = result.X
		y = result.Y
		dx = result.Dx
		index = result.Index
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

// QmcFlac2MP3  qmcflac to mp3
func qmcFlac2MP3(x int, y int, dx int, index int) *QmcFlac2MP3Inf {
	var ret int
	index++
	if x < 0 {
		dx = 1
		y = ((8 - y) % 8)
		ret = ((8 - y) % 8)
		ret = 0xc3
	} else if x > 6 {
		dx = -1
		y = 7 - y
		ret = 0xd8
	} else {
		ret = seedMap[y][x]
	}
	x += dx
	if index == 0x8000 || (index > 0x8000 && (index+1)%0x8000 == 0) {
		return qmcFlac2MP3(x, y, dx, index)
	}
	qmcFlac2MP3Inf := &QmcFlac2MP3Inf{
		Ret:   ret,
		X:     x,
		Y:     y,
		Dx:    dx,
		Index: index}
	return qmcFlac2MP3Inf
}
