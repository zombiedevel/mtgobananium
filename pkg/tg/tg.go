package tg

import (
	"github.com/zombiedevel/go-tdlib"
	"go.uber.org/zap"
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


func IsAdmin(chatID int64, userID int32 , client *tdlib.Client, log *zap.Logger) bool {
	u, err := client.GetChatMember(chatID, userID)
	if err != nil {
		log.Error("Error GetChatMember", zap.Error(err))
		return false
	}

	if u.Status.GetChatMemberStatusEnum() == tdlib.ChatMemberStatusAdministratorType || u.Status.GetChatMemberStatusEnum() == tdlib.ChatMemberStatusCreatorType {
		return true
	}

	return false
}


func IsPrivate(chatID int64, client *tdlib.Client, log *zap.Logger) bool {
	c, err := client.GetChat(chatID)
	if err != nil {
		log.Error("Error GetChat", zap.Error(err))
		return false
	}
	if c.Type.GetChatTypeEnum() == tdlib.ChatTypePrivateType {
		return true
	}
	return false
}

func CheckCommand(msgText string, entity []tdlib.TextEntity) string {
	if msgText != "" {
		if msgText[0] == '/' {
			if len(entity) >= 1 {
				if entity[0].Type.GetTextEntityTypeEnum() == "textEntityTypeBotCommand" {
					if i := strings.Index(msgText[:entity[0].Length], "@"); i != -1 {
						return msgText[:i]
					}
					return msgText[:entity[0].Length]
				}
			}
			if len(msgText) > 1 {
				if i := strings.Index(msgText, "@"); i != -1 {
					return msgText[:i]
				}
				if i := strings.Index(msgText, " "); i != -1 {
					return msgText[:i]
				}
				return msgText
			}
		}
	}
	return ""
}

func CommandArgument(msgText string) string {
	if msgText[0] == '/' {
		if i := strings.Index(msgText, " "); i != -1 {
			return msgText[i+1:]
		}
	}
	return ""
}

func GetSenderID(sender tdlib.MessageSender) int64 {
	if sender.GetMessageSenderEnum() == "messageSenderUser" {
		return int64(sender.(*tdlib.MessageSenderUser).UserID)
	}
	return sender.(*tdlib.MessageSenderChat).ChatID
}