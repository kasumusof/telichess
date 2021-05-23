package main

import (
	"fmt"
	"log"

	"github.com/kasumusof/chessgo"
	chess "github.com/kasumusof/chessgo"
	"github.com/kasumusof/telichess/models"
)

var (
	updateChan      = make(chan *Update)
	callbackChan    = make(chan string)
	repliedUpdates  = make(map[string]struct{})
	repliedMessages = make(map[int]struct{})
)

func handleUpdate(updateChan chan *Update) {
	var update *Update
	var message *Message
	var callback *CallbackQuery

	for {
		select {
		case update = <-updateChan:
			log.Println("update received")
			message = update.Message
			callback = update.CallbackQuery
		}
		if message != nil {
			handleMessage(message)
		}
		if callback != nil {
			handleCallback(callback)
		}
	}
}

func handleMessage(msg *Message) {
	chatID := msg.Chat.ID
	addUser(chatID)
	if _, ok := repliedMessages[update.Message.MessageID]; !ok {
		txt := update.Message.Text

		log.Println("Message Recieved", txt)
		if update.Message.Entities != nil {
			if update.Message.Entities[0].Type == "bot_command" {
				if txt == "/start" {
					handleStartCommand(txt)
				}
			}

		}
	} else {
		// msg, err =/ sendMessage(myBot, update.Message.Chat.ID, "Hello, please send a supported command /playgame not working right now", nil)
		if err != nil {
			log.Println("Unsupported Commnad sent:", update.Message.Text, err)
		}
	}
	repliedMessages[update.Message.MessageID] = struct{}{}
}

func handleCallback(callbackQ *CallbackQuery) {
	chatID := callbackQ.Message.Chat.ID
	addUser(chatID)
	if _, ok := repliedUpdates[callbackQ.ID]; !ok {

		log.Println("handleCallback Hit", callbackQ.Data, "with update ID", callbackQ.ID)
		callbackAction(callbackQ)

		repliedUpdates[callbackQ.ID] = struct{}{}
	}
}

func callbackAction(callbackQ *CallbackQuery) {
	var fen, text, history string
	callbackDt := callbackQ.Data
	chatID := callbackQ.Message.Chat.ID
	userID := fmt.Sprintf("%d", chatID)
	messageID := callbackQ.Message.MessageID
	// msgTxt := callbackQ.Message.Text

	_, err := answerCallbackQuery(myBot, callbackQ.ID, callbackDt, 0)
	if err != nil {
		log.Println("Error encountered: async", err)
	}
	user, err := models.GetUser(userID)

	switch callbackDt {
	case "playnow", "newgame":
		fen = chess.StartPos
		history = ""
		models.UpdateUserFEN(userID, fen, trimSpaces(history))
		text = fmt.Sprintf("%s to move.", "")
		editMessageText(myBot, chatID, messageID, "", text, genInlineKeyboardBoard(fen))
	case "undo":
		// log.Println("in undo")
		// fen = user.FEN
		history = user.History
		chessB := chess.NewStdBoard()
		dum := stringToSlice(history)
		chessB.PlayoutMoves(dum)
		if len(dum) > 1 {
			// log.Println(dum)
			chessB.UnMove()
			text = fmt.Sprintf("%s to move.", colors[int(chessB.TurnToMove())])
			fen = chessB.ToFEN()
			dum = dum[:len(dum)-1]
			history = sliceToString(dum)
			// log.Println(fen)
			models.UpdateUserFEN(userID, fen, history)
		}
		if len(dum) <= 1 {
			models.UpdateUserFEN(userID, chessgo.StartPos, trimSpaces(history))
		}

		editMessageText(myBot, chatID, messageID, "", text, genIngameInlineKeys(fen))
		log.Println("finished undo")
	default:

		if err != nil {
			log.Fatal("Error encountered", err)
		}

		chessB := chess.NewStdBoard()
		history = user.History
		// log.Println(history)
		fen = chessgo.StartPos
		if history != "" {
			dum := stringToSlice(history)
			// log.Println(dum, len(dum))
			chessB.PlayoutMoves(dum)
			fen = chessB.ToFEN()
		}
		var toMove string

		rank, file := int(chess.SquareNameToInt(callbackDt))/8, int(chess.SquareNameToInt(callbackDt))%8
		keyboard := callbackQ.Message.ReplyMarkup.InlineKeyboard
		first := keyboard[:8]
		// last := keyboard[len(keyboard)-1]

		callBackText := first[7-rank][file].Text
		log.Println("callback txt", callBackText)
		switch callBackText {
		case " ":
			editMessageText(myBot, chatID, messageID, "", text, genIngameInlineKeys(fen))
		case ".":
			// fen = chess.StartPos
			keyboard = genIngameInlineKeys(fen)
			toMove = colors[int(chessB.TurnToMove())]
			// text = fmt.Sprintf("%s to move. %s clicked on %s", toMove, "", callbackDt)

			d := callbackQ.Message.Text
			mov := fmt.Sprintf("%s%s", string(d[(len(d)-2):]), callbackDt)
			log.Println("mov", mov)
			//TODO:
			log.Println("Executed", chessB.Move(mov))
			fen = chessB.ToFEN()
			history += fmt.Sprintf(" %s", mov)
			keyboard = genIngameInlineKeys(fen)

			models.UpdateUserFEN(userID, fen, trimSpaces(history))

			text = fmt.Sprintf("%s to move. %s clicked on %s", toMove, "", callbackDt)

			editMessageText(myBot, chatID, messageID, "", text, keyboard)

		default:
			keyboard = genIngameInlineKeys(fen)
			first = keyboard[:8]
			// TODO: a := chessB.AvailableMoves(chessB.TurnToMove())
			a := chessB.GenMove(chessB.TurnToMove())
			b := filterMoves(callbackDt, a)
			log.Println("in here in moves debug", a, callbackDt)

			var pos []int
			for _, str := range b {
				val := chess.SquareNameToInt(string(str))
				pos = append(pos, val)
			}

			for _, sqr := range pos {
				var row, col int
				row, col = sqr/8, sqr%8
				first[7-row][col].Text = "."
			}
			toMove = colors[int(chessB.TurnToMove())]
			text = fmt.Sprintf("%s to move. %s clicked on %s", toMove, "", callbackDt)
			_, err = editMessageText(myBot, chatID, messageID, callbackQ.InlineMessageID, text, keyboard)
			if err != nil {
				log.Println("Error enounterd moves", err)
			}
		}

	}

}

func handleStartCommand(txt string) {
	_, err := sendMessage(
		myBot,
		update.Message.Chat.ID,
		fmt.Sprintf("Hello, I am %s. Click on PlayNow to play a game with me.", myBot.Name),
		genStartInlinKeys(),
	)
	if err != nil {
		log.Println("error in sending reply:", txt, err)
	}
}
