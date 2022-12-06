package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// bot instance
var bot *botApi.BotAPI

// forward messages to
var forwardMessagesTo int64
var confirmReceive bool = false

func initializeBot(token string) {
	//token := os.Getenv("BOT_TOKEN")
	BotAPI, err := botApi.NewBotAPI(token)

	if err != nil {
		log.Panic(err)
	}

	bot = BotAPI
}

func main() {

	str := flag.String("token", os.Getenv("BOT_TOKEN"), "Telegram bot token")
	forwardMessagesToStr := flag.String("chatid", os.Getenv("CHAT_ID"), "Chat Id to forward messages to")
	confirmReceive = *flag.Bool("confirmreceive", false, "Confirm message received to each user")
	flag.Parse()

	forwardMessagesTo = ParseInt(*forwardMessagesToStr)

	log.Printf("Messages will be forwarded to: %d", forwardMessagesTo)

	token := *str

	initializeBot(token)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := botApi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if update.Message.Chat.ID == forwardMessagesTo {
			// todo: handle replies
			continue
		}

		forwardedMsg := botApi.NewForward(forwardMessagesTo, update.Message.From.ID, update.Message.MessageID)
		if _, err := bot.Send(forwardedMsg); err != nil {
			LogToChat(err.Error())
		}

		if confirmReceive {
			SendReply(update.Message.Chat.ID, "Message received", update.Message.MessageID)
		}
	}
}

func SendReply(chatId int64, text string, replyTo ...int) {
	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	reply := botApi.NewMessage(chatId, text)
	if len(replyTo) > 0 {
		reply.ReplyToMessageID = replyTo[0]
	}
	if _, err := bot.Send(reply); err != nil {
		LogToChat(err.Error())
	}
}

func ParseInt(val string) int64 {
	chatid, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Panic(err)
	}
	return chatid
}

func LogToChat(msg string) {
	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	message := botApi.NewMessage(forwardMessagesTo, msg)
	if _, err := bot.Send(message); err != nil {
		log.Print(err)
	}
}
