package tg

import (
	"github.com/zombiedevel/go-tdlib"
	"regexp"
	"strings"
)

func TryExtractText(msg *tdlib.Message) string {
	text, ok := msg.Content.(*tdlib.MessageText)
	if ok {
		if text.Text != nil {
			return regexp.MustCompile(`(?m)<[^>]+>`).ReplaceAllLiteralString(text.Text.Text, "")
		}
	}
	return ""
}

func TryTextWithoutCommand(msg *tdlib.Message) string {
	text := TryExtractText(msg)
	commandEndPos := strings.Index(text, " ")
	if commandEndPos == -1 {
		if strings.HasPrefix(text, "/") {
			return ""
		}
		return text
	}
	return strings.Trim(text[commandEndPos:], " ")
}

func SendTextMessage(text string, chatID int64, client *tdlib.Client, markup tdlib.ReplyMarkup) {
	inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(text, nil), true, true)
	client.SendMessage(chatID, int64(0), int64(0), nil, markup, inputMsgTxt)
}


