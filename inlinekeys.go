package main

import (
	chess "github.com/kasumusof/chessgo"
)

var lichessBoard *chess.Board

func genInlineKeyboardBoard(fen string) [][]InlineKeyboardButton {
	ret := [][]InlineKeyboardButton{}
	lichessBoard = chess.NewBoardFromFEN(fen)
	for i := 0; i < 8; i++ {
		iLKey := []InlineKeyboardButton{}
		for j := 0; j < 8; j++ {
			key := InlineKeyboardButton{}
			iLKey = append(iLKey, key)
		}
		ret = append(ret, iLKey)
	}
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			text := string(chess.GetDisplay(lichessBoard.GetPiece12(i*8 + j)))
			callBackData := chess.SquareName(i*8 + j)
			// log.Println("squarename", callBackData)
			ret[7-i][j].Text = text
			ret[7-i][j].CallbackData = callBackData
		}
	}
	return ret
}

func genStartInlinKeys() [][]InlineKeyboardButton {
	ret := [][]InlineKeyboardButton{
		{
			{
				Text:         "playnow",
				CallbackData: "playnow",
			},
		},
	}
	return ret
}

func genIngameInlineKeys(fen string) [][]InlineKeyboardButton {
	ret := genInlineKeyboardBoard(fen)
	c := []InlineKeyboardButton{
		{
			Text:         "undo",
			CallbackData: "undo",
		}, {
			Text:         "newgame",
			CallbackData: "newgame",
		},
	}
	ret = append(ret, c)
	return ret
}

// // func genInlineKeyboardNewGame() [][]InlineKeyboardButton {
// // 	ret := [][]InlineKeyboardButton{
// // 		[]InlineKeyboardButton{
// // 			{
// // 				Text:         "newgame",
// // 				CallbackData: "newgame",
// // 			},
// // 		},
// // 		[]InlineKeyboardButton{
// // 			{
// // 				Text:         "continue",
// // 				CallbackData: "continue",
// // 			},
// // 		},
// // 	}
// // 	return ret
// // }

// func genMovesFromFEN(fen string) []int {
// 	res := []int{}
// 	board := chess.NewBoardFromFEN(fen)
// 	mvs := board.AvailableMoves(board.TurnToMove())
// 	fmt.Println(mvs)
// 	return res
// }

// func filterMoves(mv string, mvs []string) []string {
// 	res := []string{}
// 	if mv == "" {
// 		return res
// 	}
// 	for _, str := range mvs {
// 		if mv == string(str[0:2]) {
// 			res = append(res, str[2:])
// 		}
// 	}
// 	return res
// }
