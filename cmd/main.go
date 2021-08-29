package main

import (
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/zombiedevel/go-tdlib"
	"github.com/zombiedevel/mtgobabanium/handlers"
	"github.com/zombiedevel/mtgobabanium/pkg/tg"
	"go.uber.org/zap"
)


func main() {
	_ = godotenv.Load(".env")

	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	var config map[string]string
	config, err := godotenv.Read()
	if err != nil {
		logger.Error("Error load env", zap.Error(err))
		return
	}
	// Create new instance of client
	client := tdlib.NewClient(tdlib.Config{
		APIID:               config["APP_ID"],
		APIHash:             config["APP_HASH"],
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./tdlib-db",
		FileDirectory:       "./tdlib-files",
		IgnoreFileNames:     false,
	})

	go func() {
		eventFilter := func(msg *tdlib.TdMessage) bool {
			return true
		}
		receiver := client.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 1000)
		for newMsg := range receiver.Chan {
			update := (newMsg).(*tdlib.UpdateNewMessage)
			if update.Message.ReplyToMessageID > 0 {
				msgData, err := client.GetMessage(update.Message.ChatID, update.Message.ReplyToMessageID)
				if err != nil {
					logger.Error("get message reply", zap.Error(err))
				}
				member, err := client.GetChatMember(update.Message.ChatID, update.Message.Sender.(*tdlib.MessageSenderUser).UserID)
				if err != nil {
					logger.Error("Error get member", zap.Error(err))
				}
				cmd := tg.TryExtractText(update.Message)
				switch cmd {
				case "!src": handlers.SrcHandler(msgData, client, logger)
				case "!ban":
					if member.Status.GetChatMemberStatusEnum() == tdlib.ChatMemberStatusAdministratorType {
						// TODO: Make ban handler
					}
				case "!ro":
					if member.Status.GetChatMemberStatusEnum() == tdlib.ChatMemberStatusAdministratorType {
						handlers.RoHandler(msgData, client, logger)
					}
				case "!report": handlers.ReportHandler(msgData, client, logger)
				case "О боте": // TODO: Make about handler
				}
			}
			switch tg.TryExtractText(update.Message) {
			case "/start":
				handlers.StartHandler(newMsg, client, logger)
			}
		}
	}()
	for {
		currentState, _ := client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			_, err := client.CheckAuthenticationBotToken(config["BOT_TOKEN"])
			if err != nil {
				logger.Error("Error check bot token", zap.Error(err))
				return
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			user, err := client.GetMe()
			if err != nil {
				logger.Error("Error client GetMe", zap.Error(err))
			}
			logger.Info("Authorization bot success", zap.String("username", user.Username), zap.String("FirstName", user.FirstName))
			break
		}
	}

	// rawUpdates gets all updates comming from tdlib
	rawUpdates := client.GetRawUpdatesChannel(100)
	for range rawUpdates {
		// Show all updates
		//fmt.Println(upd.Data)
		//message := update.Data
	}
}
