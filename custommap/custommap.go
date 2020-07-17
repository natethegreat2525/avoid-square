package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/natethegreat2525/avoid-square/utils"
)

const maxSize = 1600000000 // 1.6 Billion x8 bytes ~ 10 Gb

var checked *[maxSize]uint64

var numSlots = uint64(len(checked))
var numChecked = 0

var checkedHigh = make(map[int64]int8)
var cnt = 0

func setCheck(bh int64, val int8, movesLeft byte) {
	if movesLeft <= 16 {
		id := utils.Qhash(uint64(bh)) % numSlots

		if (checked[id] & 0xff00000000000000) == 0 {
			numChecked++
		}
		checked[id] = uint64(bh) | (uint64(val+2) << (7 * 8))
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
			if utils.CheckBoard(newBoard) {
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
			if !utils.CheckBoard(newBoard) {
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
			if utils.CheckBoard(newBoard) && heurBoard[i] < lowestH {
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
	bh := utils.ByteToInt64(utils.Hash(board))
	if movesLeft <= 16 {
		id := utils.Qhash(uint64(bh)) % numSlots
		cid := checked[id]
		cbh := cid & 0x00ffffffffffffff
		cval := (cid >> (7 * 8)) & 0x00000000000000ff
		if cval != 0 && cbh == uint64(bh) {
			return int8(cval) - 2
		}
	} else {
		if val, ok := checkedHigh[bh]; ok {
			return val
		}
	}
	curCheck := utils.CheckBoard(board)
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
		fmt.Println(numChecked, len(checkedHigh))
	}

	score := worstScore
	if turn == 1 {
		score = bestScore
	}

	setCheck(bh, score, movesLeft)
	return score
}

func main() {
	// Must do this or go will double memory usage!
	debug.SetGCPercent(5)
	checked = &[maxSize]uint64{}

	for i := 1; i < 25; i++ {
		testBoard := []byte{
			1, 0, 0, 0, 0,
			0, 0, 0, 0, 0,
			0, 0, 0, 0, 0,
			0, 0, 0, 0, 0,
			0, 0, 0, 0, 0,
		}
		testBoard[i] = 2
		fmt.Printf("IDX: %v\nScore: %v\n", i, scorePosition(testBoard, 1, 24))
	}

	f, err := os.Create("./out_boards_upper_left")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer f.Close()

	for board, score := range checkedHigh {
		boardscore := []byte{utils.GetB(board, 0), utils.GetB(board, 1), utils.GetB(board, 2), utils.GetB(board, 3), utils.GetB(board, 4), utils.GetB(board, 5), utils.GetB(board, 6), byte(score)}
		f.Write(boardscore)
	}
}
