package utils

import (
	"fmt"
	"math"
)

//BoardToString easy board printing
func BoardToString(board []byte) string {
	length := int(math.Sqrt(float64(len(board))))
	s := ""
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			s = fmt.Sprintf("%v%v", s, board[i+j*length])
		}
		s = s + "\n"
	}
	return s
}

//UintToBoard converts a uint to a non compacted byte slice board
func UintToBoard(enc uint64, size int) []byte {
	board := make([]byte, size)
	for i := 0; i < 25; i++ {
		val := enc & 0x0000000000000003
		board[i] = byte(val)
		enc = enc >> 2
	}
	return board
}

//BoardToEnc bit packs board into 7 bytes
func BoardToEnc(board []byte) []byte {
	s := make([]byte, 0)
	var acc byte
	var curAc byte
	for i := 0; i < len(board); i++ {
		if acc == 4 {
			acc = 0
			s = append(s, curAc)
			curAc = 0
		}
		curAc = curAc | (board[i] << uint(2*acc))
		acc++
	}
	s = append(s, curAc)
	return s
}

//CheckBoard checks if a board is a valid (non losing) position for both players
func CheckBoard(board []byte) bool {
	length := int(math.Sqrt(float64(len(board))))
	for i := 0; i < length-1; i++ {
		for j := 0; j < length-1; j++ {
			for x := 1; x < length-i; x++ {
				for y := 0; y < length-j; y++ {
					typ := board[i+j*length]
					if typ == 0 {
						continue
					}
					found := false
					xx := x
					yy := y
					ii := i + x
					jj := j + y
					for n := 0; n < 3; n++ {
						if ii < 0 || jj < 0 || ii >= length || jj >= length {
							found = true
							break
						}
						if board[ii+jj*length] != typ {
							found = true
							break
						}
						tmp := xx
						xx = -yy
						yy = tmp
						ii = ii + xx
						jj = jj + yy
					}
					if !found {
						return false
					}
				}
			}
		}
	}
	return true
}

//Hash hashes board to byte slice accounting for all symmetries
func Hash(board []byte) []byte {
	c12x := 0.0
	c12y := 0.0
	length := int(math.Sqrt(float64(len(board))))
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			if board[i+j*length] == 1 {
				c12x = c12x + float64(i) + .5 - float64(length)/2.0
				c12y = c12y + float64(j) + .5 - float64(length)/2.0
			}
			if board[i+j*length] == 2 {
				c12x = c12x + .9*(float64(i)+.5-float64(length)/2.0)
				c12y = c12y + .9*(float64(j)+.5-float64(length)/2.0)
			}
		}
	}
	flip := math.Abs(c12x) > math.Abs(c12y)
	newBoard := make([]byte, len(board))
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			x := i
			y := j
			if c12x < 0 {
				x = length - 1 - x
			}
			if c12y < 0 {
				y = length - 1 - y
			}
			a := i
			b := j
			if flip {
				t := a
				a = b
				b = t
			}
			newBoard[a+b*length] = board[x+y*length]
		}
	}
	return BoardToEnc(newBoard)
}

//ByteToInt64 convert 8 or fewer bytes to int64
func ByteToInt64(val []byte) int64 {
	var ret int64
	for i := 0; i < len(val); i++ {
		ret = ret | (int64(val[i]) << uint(8*i))
	}
	return ret
}

//Qhash quick hash function from stack overflow
func Qhash(x uint64) uint64 {
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	x = (x ^ (x >> 31))
	return x
}

//GetB get indexed byte from int64 (slow)
func GetB(v int64, idx int) byte {
	return byte((v >> (8 * idx)) & 0x00000000000000ff)
}
