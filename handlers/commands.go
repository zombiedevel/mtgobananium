package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/zombiedevel/go-tdlib"
	"github.com/zombiedevel/mtgobabanium/pkg/tg"
	"go.uber.org/zap"
	"time"
)

func StartHandler(msg tdlib.TdMessage, client *tdlib.Client, log *zap.Logger) {
	upd := (msg).(*tdlib.UpdateNewMessage)
	sender := upd.Message.Sender.(*tdlib.MessageSenderUser)
	user, err := client.GetUser(int32(sender.UserID))

	if err != nil {
		log.Error("Error get user", zap.Error(err))
	}
	log.Info("StartHandler",
		zap.Int32("ID", user.ID),
		zap.String("Username", user.Username),
		zap.String("FirstName", user.FirstName),
		zap.String("LastName", user.LastName),
	)
	var buttons [][]tdlib.KeyboardButton
	buttons = append(buttons, []tdlib.KeyboardButton{
		*tdlib.NewKeyboardButton("Пройти проверку", tdlib.NewKeyboardButtonTypeText()),
		*tdlib.NewKeyboardButton("О боте", tdlib.NewKeyboardButtonTypeText()),
	},
	)
	var format *tdlib.FormattedText
	format = tdlib.NewFormattedText(fmt.Sprintf("Привет, %s", user.FirstName), nil)
	text := tdlib.NewInputMessageText(format, false, false)
	client.SendMessage(upd.Message.ChatID, 0,
		0,
		tdlib.NewMessageSendOptions(false, true, nil),
		tdlib.ReplyMarkup(tdlib.NewReplyMarkupShowKeyboard(buttons, true, false, true)), text)
	return
}
// Restrict member handler
func RoHandler(msg *tdlib.Message, client *tdlib.Client, log *zap.Logger) {
	t := time.Now()
	uTime := t.Local().Add(time.Minute * 15).Unix()
	userId := msg.Sender.(*tdlib.MessageSenderUser).UserID
	user, err := client.GetUser(userId)
	if err != nil {
		log.Error("Error GetUser", zap.Error(err))
		return
	}
	if err := restrict(msg.Sender.(*tdlib.MessageSenderUser), msg.ChatID, client, uTime); err != nil {
		log.Error("Error restrict user", zap.Error(err))
	}
	tg.SendTextMessage(fmt.Sprintf("Пользователь %s помещён в карантин.", user.FirstName), msg.ChatID, client, nil)
	return
}

func SrcHandler(msg *tdlib.Message, client *tdlib.Client, log *zap.Logger) {
	byte, err := json.MarshalIndent(msg, "", " ")
	if err != nil {
		log.Error("Error MarshalIndent", zap.Error(err))
		return
	}
	var format *tdlib.FormattedText
	format, err = client.ParseTextEntities(fmt.Sprintf("```%s```", string(byte)), tdlib.NewTextParseModeMarkdown(2))
	if err != nil {
		log.Error("Error ParseTextEntities", zap.Error(err))
	}
	inputMsgTxt := tdlib.NewInputMessageText(format, true, false)

	client.SendMessage(msg.ChatID, msg.MessageThreadID, msg.ID, nil, nil, inputMsgTxt)
}

func BanHandler(msg *tdlib.Message, client *tdlib.Client, log *zap.Logger) {
	userId := msg.Sender.(*tdlib.MessageSenderUser).UserID
	user, err := client.GetUser(userId)
	if err != nil {
		log.Error("Error GetUser", zap.Error(err))
		return
	}
	if _, err := client.SetChatMemberStatus(msg.ChatID, user.ID, tdlib.NewChatMemberStatusBanned(0)); err != nil {
		log.Error("Error SetChatMemberStatus", zap.Error(err))
		return
	}
	if _, err := client.DeleteChatMessagesFromUser(msg.ChatID, user.ID); err != nil {
		log.Error("Error DeleteChatMessagesFromUser", zap.Error(err))
		return
	}
	msgText := fmt.Sprintf("Пользователь %s утилизирован.", user.FirstName)
	tg.SendTextMessage(msgText, msg.ChatID, client, nil)
}

func ReportHandler(msg *tdlib.Message, client *tdlib.Client, log *zap.Logger) {

}