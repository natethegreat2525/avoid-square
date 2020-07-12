package main

import (
	"fmt"
	"math"
)

func boardToString(board []byte) string {
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

func boardToEnc(board []byte) []byte {
	s := make([]byte, 0)
	var acc byte
	var curAc byte
	for i := 0; i < len(board); i++ {
		if acc == 4 {
			acc = 0
			s = append(s, curAc)
			curAc = 0
		}
		curAc = curAc | (board[i] << (2 * acc))
		acc++
	}
	s = append(s, curAc)
	return s
}

func checkBoard(board []byte) bool {
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

func hash(board []byte) []byte {
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
	return boardToEnc(newBoard)
}

func byteToInt64(val []byte) int64 {
	var ret int64 = 0
	for i := 0; i < len(val); i++ {
		ret = ret | (int64(val[i]) << (8 * i))
	}
	return ret
}

var checked map[int64]int8 = make(map[int64]int8)
var checkedHigh map[int64]int8 = make(map[int64]int8)
var cnt = 0

func setCheck(bh int64, val int8, movesLeft byte) {
	if movesLeft <= 16 {
		checked[bh] = int8(val)
		if len(checked) > 100000000 {
			checked = make(map[int64]int8)
		}
	} else {
		checkedHigh[bh] = int8(val)
	}
}

func aiMove(board []byte, turn byte) ([]byte, bool) {
	moveCount := 0
	for i := 0; i < len(board); i++ {
		if board[i] == 0 {
			moveCount++
			newBoard := make([]byte, len(board))
			copy(newBoard, board)
			newBoard[i] = turn
			if checkBoard(newBoard) {
				return newBoard, true
			}
		}
	}
	if moveCount == 0 {
		return nil, false
	}

	for i := 0; i < len(board); i++ {
		if board[i] == 0 {
			newBoard := make([]byte, len(board))
			copy(newBoard, board)
			newBoard[i] = turn
			return newBoard, true
		}
	}

	panic("not popsicle")
}

func aiSmart(board []byte, turn byte) ([]byte, bool) {
	heurBoard := make([]int, len(board))
	for i := 0; i < len(board); i++ {
		if board[i] == 0 {
			newBoard := make([]byte, len(board))
			copy(newBoard, board)
			newBoard[i] = 3 - turn
			// if opponent fails by moving here
			// then increment that spot, don't move there
			if !checkBoard(newBoard) {
				heurBoard[i]++
			}
		}
	}

	moveCount := 0
	lowestH := 99
	var lowestBoard []byte
	for i := 0; i < len(board); i++ {
		if board[i] == 0 {
			moveCount++
			newBoard := make([]byte, len(board))
			copy(newBoard, board)
			newBoard[i] = turn
			if checkBoard(newBoard) && heurBoard[i] < lowestH {
				lowestH = heurBoard[i]
				lowestBoard = newBoard
			}
		}
	}
	// if low heuristic is found, return the board that has the lowest heuristic
	if lowestH < 99 {
		return lowestBoard, true
	}
	//if board is full, cant move
	if moveCount == 0 {
		return nil, false
	}

	//just move the first move if we couldnt find a good move
	for i := 0; i < len(board); i++ {
		if board[i] == 0 {
			newBoard := make([]byte, len(board))
			copy(newBoard, board)
			newBoard[i] = turn
			return newBoard, true
		}
	}

	panic("not popsicle")
}

func scorePosition(board []byte, turn byte, movesLeft byte) int8 {
	bh := byteToInt64(hash(board))
	if movesLeft <= 16 {
		if val, ok := checked[bh]; ok {
			return val
		}
	} else {
		if val, ok := checkedHigh[bh]; ok {
			return val
		}
	}
	curCheck := checkBoard(board)
	if !curCheck {
		v := -1
		if turn == 1 {
			v = 1
		}
		setCheck(bh, int8(v), movesLeft)
		return int8(v)
	}

	bestScore := int8(-9)
	worstScore := int8(9)
	moves := 0

	nextTurn := 3 - turn
	//if turn == 1 || movesLeft > 8 {
	for i := 0; i < len(board); i++ {
		if board[i] == 0 {
			moves++
			newBoard := make([]byte, len(board))
			copy(newBoard, board)
			newBoard[i] = turn
			sc := scorePosition(newBoard, nextTurn, movesLeft-1)
			if bestScore < sc {
				bestScore = sc
			}
			if worstScore > sc {
				worstScore = sc
			}
			if turn == 1 && bestScore > 0 {
				break
			}
			if turn != 1 && worstScore < 0 {
				break
			}
		}
	}
	/*} else {
		newBoard, valid := aiSmart(board, turn)
		if valid {
			moves++
			sc := scorePosition(newBoard, nextTurn, movesLeft-1)
			if bestScore < sc {
				bestScore = sc
			}
			if worstScore > sc {
				worstScore = sc
			}
		}
	}*/

	if moves == 0 {
		setCheck(bh, 0, movesLeft)
		return 0
	}
	cnt++
	if cnt >= 100000 {
		cnt = 0
		fmt.Println(len(checked), len(checkedHigh))
	}

	score := worstScore
	if turn == 1 {
		score = bestScore
	}

	setCheck(bh, score, movesLeft)
	return score
}
func main() {
	fmt.Printf("%v\n", scorePosition([]byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}, 1, 25))
}

/*
func main() {
	db, _ := bitcask.Open("./testdb")
	defer db.Close()

	for i := 0; i < 1000000; i++ {
		db.Put([]byte(fmt.Sprintf("%v", i)), []byte(fmt.Sprintf("test %v", i)))
	}
	//val, _ := db.Get([]byte("Hello"))
	//fmt.Printf("%v\n", val)
	fmt.Println("made db")
	// for i := 0; i < 1000000; i++ {
	// 	db.Get([]byte(fmt.Sprintf("%v", i)))
	// }
}*/
