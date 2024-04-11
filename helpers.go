package main

import (
	"flag"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	"strconv"
	"strings"
)

func createMedia(media tele.Media, caption ...string) tele.Inputtable {
	mediaType := media.MediaType()
	file := *media.MediaFile()
	return createMediaFromType(mediaType, file, caption...)
}

func createMediaFromType(mediaType string, file tele.File, caption ...string) tele.Inputtable {
	var c = ""
	if len(caption) > 0 {
		c = strings.Join(caption, "")
	}
	switch mediaType {
	case "video":
		return &tele.Video{File: file, Caption: c}
	case "photo":
		return &tele.Photo{File: file, Caption: c}
	case "audio":
		return &tele.Audio{File: file, Caption: c}
	default:
		logToChat(fmt.Sprintf("Type %s not supported for album", mediaType))
	}
	return nil
}

func parseConfigParams() ConfigParams {
	str := flag.String("token", os.Getenv("BOT_TOKEN"), "Telegram bot token")
	forwardMessagesToStr := flag.String("chatid", os.Getenv("CHAT_ID"), "Chat Id to forward messages to")
	confirmReceive := flag.Bool("confirmreceive", false, "Confirm message received to each user")
	useLogger := flag.Bool("uselogger", false, "Use verbose logger")
	flag.Parse()

	forwardMessagesTo := parseInt(*forwardMessagesToStr)
	config := ConfigParams{
		botToken:          *str,
		forwardMessagesTo: forwardMessagesTo,
		confirmReceive:    *confirmReceive,
		useLogger:         *useLogger,
	}
	return config
}

func parseInt(val string) int64 {
	chatid, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Panic(err)
	}
	return chatid
}

func logToChat(msg string) {
	if _, err := b.Send(forwardTo, msg); err != nil {
		log.Print(err)
	}
}
