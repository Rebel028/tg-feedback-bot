package main

import (
	"flag"
	"fmt"
	"gopkg.in/telebot.v3/middleware"
	"log"
	"os"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

var b *tele.Bot

var forwardTo *tele.Chat
var confirmReceive bool = false

var mediaGroups = make(map[string][]tele.Media, 3)

type ConfigParams struct {
	botToken          string
	forwardMessagesTo int64
	confirmReceive    bool
	useLogger         bool
}

func initializeBot(token string) *tele.Bot {

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	return bot
}

func main() {

	configParams := ParseConfigParams()

	forwardTo = &tele.Chat{ID: configParams.forwardMessagesTo}
	confirmReceive = configParams.confirmReceive

	b = initializeBot(configParams.botToken)
	log.Printf("Messages will be forwarded to: %d", configParams.forwardMessagesTo)

	if configParams.useLogger {
		b.Use(middleware.Logger())
	}

	//default handlers
	b.Handle(tele.OnText, ForwardMessage)
	b.Handle(tele.OnVoice, ForwardMessage)
	b.Handle(tele.OnVideoNote, ForwardMessage)
	b.Handle(tele.OnDocument, ForwardMessage)
	b.Handle(tele.OnAnimation, ForwardMessage)
	b.Handle(tele.OnLocation, ForwardMessage)

	//special handler for albums
	b.Handle(tele.OnMedia, CheckForwardAlbum)

	b.Start()
}

func CheckForwardAlbum(c tele.Context) error {
	mediaGroupId := c.Message().AlbumID
	if mediaGroupId != "" {
		HandleAlbum(b, c, mediaGroupId)
		return nil
	}
	return ForwardMessage(c)
}

func ForwardMessage(c tele.Context) error {
	if _, err := b.Forward(forwardTo, c.Message()); err != nil {
		LogToChat(err.Error())
		return err
	}
	if confirmReceive {
		err := c.Reply("Message received!")
		if err != nil {
			LogToChat(err.Error())
			return err
		}
	}
	return nil
}

func HandleAlbum(b *tele.Bot, c tele.Context, mediaGroupId string) {
	media := c.Message().Media()
	log.Printf("Adding file to arr: %s", media.MediaFile().FileID)

	if messages, ok := mediaGroups[mediaGroupId]; ok {
		mediaGroups[mediaGroupId] = append(messages, media)
	} else {
		mediaGroups[mediaGroupId] = []tele.Media{media}
		var caption = fmt.Sprintf("<b>Forwarded from</b> <a href=\"tg://user?id=%d\">%s</a>\n<b>Message</b>:\n %s",
			c.Message().Sender.ID,
			c.Message().Sender.FirstName+c.Message().Sender.LastName,
			c.Message().Text)

		time.AfterFunc(time.Second, func() {
			album := tele.Album{}
			for i, media := range mediaGroups[mediaGroupId] {
				fileId := media.MediaFile().FileID
				log.Printf("Sending file %s", fileId)
				var c = ""
				if i == 0 {
					c = caption
				}
				album = append(album, CreateMedia(media, c))
			}
			if _, err := b.SendAlbum(forwardTo, album, &tele.SendOptions{
				ParseMode: "Html",
			}); err != nil {
				LogToChat(err.Error())
			}
		})
	}
}

func CreateMedia(media tele.Media, caption string) tele.Inputtable {
	switch media.MediaType() {
	case "video":
		return &tele.Video{File: *media.MediaFile(), Caption: caption}
	case "photo":
		return &tele.Photo{File: *media.MediaFile(), Caption: caption}
	case "audio":
		return &tele.Audio{File: *media.MediaFile(), Caption: caption}
	default:
		LogToChat(fmt.Sprintf("Type %s not supported for album", media.MediaType()))
	}
	return nil
}

func ParseConfigParams() ConfigParams {
	str := flag.String("token", os.Getenv("BOT_TOKEN"), "Telegram bot token")
	forwardMessagesToStr := flag.String("chatid", os.Getenv("CHAT_ID"), "Chat Id to forward messages to")
	confirmReceive := flag.Bool("confirmreceive", false, "Confirm message received to each user")
	useLogger := flag.Bool("uselogger", false, "Use verbose logger")
	flag.Parse()

	forwardMessagesTo := ParseInt(*forwardMessagesToStr)
	config := ConfigParams{
		botToken:          *str,
		forwardMessagesTo: forwardMessagesTo,
		confirmReceive:    *confirmReceive,
		useLogger:         *useLogger,
	}
	return config
}

func ParseInt(val string) int64 {
	chatid, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Panic(err)
	}
	return chatid
}

func LogToChat(msg string) {
	if _, err := b.Send(forwardTo, msg); err != nil {
		log.Print(err)
	}
}
