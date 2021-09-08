package handlers

import (
	"fmt"
	"github.com/zombiedevel/go-tdlib"
	"github.com/zombiedevel/mtgobabanium/pkg/tg"
	"go.uber.org/zap"
)

func Protect(member *tdlib.ChatMember, client *tdlib.Client, chatID int64, log *zap.Logger) {
	if err := restrict(&tdlib.MessageSenderUser{UserID: member.UserID}, chatID, client, 0); err != nil {
		log.Error("Error restrict user", zap.Error(err))
		return
	}
	var buttons [][]tdlib.InlineKeyboardButton
	buttons = append(buttons, []tdlib.InlineKeyboardButton{
		*tdlib.NewInlineKeyboardButton("Пройти проверку", tdlib.NewInlineKeyboardButtonTypeURL(fmt.Sprintf("https://t.me/pornuhobot?start=protect"))),
	})
	user, err := client.GetUser(member.UserID)
	if err != nil {
		log.Error("Error GetUser", zap.Error(err))
		return
	}
	tg.SendTextMessage(fmt.Sprintf("Привет, %s! Для того что бы начать общаться в группе, тебе необходимо пройти проверку.\nЕсли ты не пройдёшь проверку в течении часа. Твой аккаунт будет заблокирован в группе.", user.FirstName), chatID, client, tdlib.NewReplyMarkupInlineKeyboard(buttons))
}

// restrict chat member
func restrict(user *tdlib.MessageSenderUser, chatID int64, client *tdlib.Client, time int64) error {
	if _, err := client.SetChatMemberStatus(
		chatID,
		user.UserID,
		tdlib.NewChatMemberStatusRestricted(true, int32(time), &tdlib.ChatPermissions{CanSendMessages: false})); err != nil {
		return err
	}
	return nil
}

func ProtectMeHandler(msg *tdlib.Message, client *tdlib.Client, log *zap.Logger) {

}