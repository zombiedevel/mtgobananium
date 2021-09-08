package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/zombiedevel/go-tdlib"
	"github.com/zombiedevel/mtgobabanium/internal/gentext"
	"github.com/zombiedevel/mtgobabanium/pkg/template"
	"github.com/zombiedevel/mtgobabanium/pkg/tg"
	"github.com/zombiedevel/mtgobabanium/pkg/tv"
	"go.uber.org/zap"
	"time"
)
func StartHandler(msg *tdlib.Message, client *tdlib.Client, log *zap.Logger) {
	user, err := client.GetUser(msg.Sender.(*tdlib.MessageSenderUser).UserID)
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
	client.SendMessage(msg.ChatID, 0,
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
	return
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
	if _, err := client.DeleteMessages(msg.ChatID, []int64{msg.ID}, true); err != nil {
		log.Error("Error DeleteMessages", zap.Error(err))
	}

	msgText := fmt.Sprintf("Пользователь %s утилизирован.", user.FirstName)
	tg.SendTextMessage(msgText, msg.ChatID, client, nil)
	return
}

func ReportHandler(msg *tdlib.Message, adminsChannelID int64,  client *tdlib.Client, log *zap.Logger) {
	channel, err := client.GetChat(adminsChannelID)
	if err != nil {
		log.Error("Error GetChat", zap.Error(err))
		return
	}

	_, err = client.ForwardMessages(channel.ID, msg.ChatID, []int64{msg.ID}, nil, false, false)
	if err != nil {
		log.Error("Error ForwardMessages", zap.Error(err))
		return
	}
	var format *tdlib.FormattedText
	format = tdlib.NewFormattedText("Мы примем все необходимые меры, спасибо.", nil)
	msgInput := tdlib.NewInputMessageText(format, true, true)

	if _, err := client.SendMessage(msg.ChatID, msg.MessageThreadID, msg.ID, nil, nil, msgInput); err != nil {
		log.Error("Error sendMessage", zap.Error(err))
		return
	}
}

func TvHandler(msg *tdlib.Message, client *tdlib.Client, log *zap.Logger) {
	moar, inputMsg := NewTvMessage(log)
	_, err := client.SendMessage(msg.ChatID, msg.MessageThreadID, msg.ID, nil, tdlib.NewReplyMarkupInlineKeyboard(moar), inputMsg)
	if err != nil {
		log.Error("Error sendMessage", zap.Error(err))
		return
	}

	return
}

func NewTvMessage(log *zap.Logger) ([][]tdlib.InlineKeyboardButton, *tdlib.InputMessageAnimation) {
	var moar [][]tdlib.InlineKeyboardButton
	moar = append(moar, []tdlib.InlineKeyboardButton{
		*tdlib.NewInlineKeyboardButton("MOAR", tdlib.NewInlineKeyboardButtonTypeCallback([]byte("next_movie"))),
	})
	video := tv.GetMovie(log)
	inputMsg := tdlib.NewInputMessageAnimation(tdlib.NewInputFileLocal(video.VideoPath), nil, nil, 0, 300, 300,  tdlib.NewFormattedText(video.Description, nil))
	return moar, inputMsg
}

func GptHandler(msg *tdlib.Message, client *tdlib.Client, log *zap.Logger) {
	message, err := client.GetMessage(msg.ChatID, msg.ReplyToMessageID)
	if err != nil {
		log.Error("Error GetMessage", zap.Error(err))
		return
	}
	if message.Content.GetMessageContentEnum() != "messageText" { return }
	replyMsgText := message.Content.(*tdlib.MessageText).Text.Text
	gpt := gentext.NewGPT3()
	gptText, err := gpt.Query(replyMsgText)
	if err != nil {
		log.Error("Error GPT3", zap.Error(err))
		return
	}
	var format *tdlib.FormattedText
	format = tdlib.NewFormattedText(gptText, nil)
	msgInput := tdlib.NewInputMessageText(format, true, true)
	if _, err := client.SendMessage(msg.ChatID, msg.MessageThreadID, message.ID, nil, nil, msgInput); err != nil {
		log.Error("Error sendMessage", zap.Error(err))
		return
	}
	if _, err := client.DeleteMessages(msg.ChatID, []int64{msg.ID}, true); err != nil {
		log.Error("Error DeleteMessages", zap.Error(err))
		return
	}
  return
}

func BioHandler(msg *tdlib.Message, client *tdlib.Client, log *zap.Logger) {
	templateStr := fmt.Sprintf(`ID: {{.ID}}
Имя: {{.FirstName}} {{.LastName}}
Имя пользователя: @{{.Username}}
О себе: {{.Bio}}`)
	member, err := client.GetChatMember(msg.ChatID, msg.Sender.(*tdlib.MessageSenderUser).UserID)
	if err != nil {
		log.Error("Error GetChatMember", zap.Error(err))
		return
	}
	user, err := client.GetUser(member.UserID)
	if err != nil {
		log.Error("Error GetUser", zap.Error(err))
		return
	}
	full, err := client.GetUserFullInfo(member.UserID)
	if err != nil {
		log.Error("Error GetUserFullInfo", zap.Error(err))
		return
	}
	text, err := template.Template("bio", templateStr, struct {
		*tdlib.User
		Bio string
	}{
		user,
		full.Bio,
	})
	if err != nil {
		log.Error("Error parsing template", zap.Error(err))
		return
	}

	var message tdlib.InputMessageContent
	message = tdlib.NewInputMessageText(tdlib.NewFormattedText(text, nil), true, true)
	if user.ProfilePhoto != nil {
		avatar, err := client.DownloadFile(user.ProfilePhoto.Big.ID, 1, 0, 0, true)
		if err != nil {
			log.Error("Error DownloadFile", zap.Error(err))
			return
		}
		if user.ProfilePhoto.HasAnimation {
			message = tdlib.NewInputMessageAnimation(
				tdlib.NewInputFileLocal(avatar.Local.Path),
				nil, nil, 0, 300, 300, tdlib.NewFormattedText(text, nil))
		} else {
			message = tdlib.NewInputMessagePhoto(
				tdlib.NewInputFileLocal(avatar.Local.Path),
				nil, nil, 300, 300, tdlib.NewFormattedText(text, nil), 0)
		}
	}
	if _, err := client.SendMessage(msg.ChatID, 0, msg.ReplyToMessageID, nil, nil, message); err != nil {
		log.Error("Error sendMessage", zap.Error(err))
		return
	}
	return
}