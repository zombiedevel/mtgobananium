package handlers

import (
	"github.com/zombiedevel/go-tdlib"
	"go.uber.org/zap"
)

func NextMovieCallback(msg *tdlib.Message, client *tdlib.Client, log *zap.Logger) {
	if msg.Content.GetMessageContentEnum() != tdlib.MessageAnimationType {
		log.Error("Message callback next_movie is not animation")
		return
	}
	moar, inputMsg := NewTvMessage(log)
   _, err := client.EditMessageMedia(msg.ChatID, msg.ID, tdlib.NewReplyMarkupInlineKeyboard(moar), inputMsg)
   if err != nil {
   	log.Error("Error EditMessageMedia", zap.Error(err))
   	return
   }
   return
}
