package main

import (
	"fmt"
	"io/ioutil"

	"github.com/natethegreat2525/avoid-square/utils"
)

var boardMap = make(map[int64]int8)

func showWinningMovesAgainst(board []byte) {
	for i := 0; i < len(board); i++ {
		if board[i] == 0 {
			newBoard := make([]byte, len(board))
			copy(newBoard, board)
			newBoard[i] = 2
			newHash := utils.ByteToInt64(utils.Hash(newBoard))
			if val, ok := boardMap[newHash]; ok && val == -1 {
				fmt.Printf("Original Board:\n%v\nWinning Move:\n%v\n", utils.BoardToString(board), utils.BoardToString(newBoard))
				break
			}
		}
	}
}

/*
func main() {
	data, err := ioutil.ReadFile("out_boards")
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}

	for j := 0; j < len(data)/8; j++ {
		i := j * 8
		boardDat := uint64(data[i]) |
			(uint64(data[i+1]) << 8) |
			(uint64(data[i+2]) << 16) |
			(uint64(data[i+3]) << 24) |
			(uint64(data[i+4]) << 32) |
			(uint64(data[i+5]) << 40) |
			(uint64(data[i+6]) << 48)
		val := int8(data[i+7])

		board := utils.UintToBoard(boardDat, 25)
		boardMap[utils.ByteToInt64(utils.Hash(board))] = val
		//fmt.Printf("Board:\n%v\nValue: %v\n\n", utils.BoardToString(board), val)
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			board := make([]byte, 25)
			board[0] = 1
			board[i+j*5] = 2
			fmt.Printf("IDX: %v Val: %v\n", i+j*5, boardMap[utils.ByteToInt64(utils.Hash(board))])
		}
	}
}
*/

func main() {
	data, err := ioutil.ReadFile("out_boards")
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}

	for j := 0; j < len(data)/8; j++ {
		i := j * 8
		boardDat := uint64(data[i]) |
			(uint64(data[i+1]) << 8) |
			(uint64(data[i+2]) << 16) |
			(uint64(data[i+3]) << 24) |
			(uint64(data[i+4]) << 32) |
			(uint64(data[i+5]) << 40) |
			(uint64(data[i+6]) << 48)
		val := int8(data[i+7])

		board := utils.UintToBoard(boardDat, 25)
		boardMap[utils.ByteToInt64(utils.Hash(board))] = val
		//fmt.Printf("Board:\n%v\nValue: %v\n\n", utils.BoardToString(board), val)
	}

	for i := 0; i < 3; i++ {
		for j := 0; j <= i; j++ {
			board := make([]byte, 25)
			board[i+j*5] = 1
			showWinningMovesAgainst(board)
		}
	}
}
