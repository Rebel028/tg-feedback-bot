package main

import (
	"github.com/glebarez/sqlite"
	"gopkg.in/telebot.v3/middleware"
	"gorm.io/gorm"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

var b *tele.Bot

// chat we want to forward messages to
var forwardTo *tele.Chat

// do we need to reply with confirmation on each message
var confirmReceive bool = false

// Telegram bots receive albums as separate messages.
// So to be able to forward them properly
// we collect all incoming media group items
// in a dictionary with media_group_id as key for each group
var mediaGroups = make(map[string][]MediaGroupItem, 3)

func main() {

	configParams := parseConfigParams()

	forwardTo = &tele.Chat{ID: configParams.forwardMessagesTo}
	confirmReceive = configParams.confirmReceive

	b = initializeBot(configParams.botToken)
	log.Printf("Messages will be forwarded to: %d", configParams.forwardMessagesTo)

	if configParams.useLogger {
		b.Use(middleware.Logger())
	}

	initializeSqlite()

	//default handlers
	b.Handle(tele.OnText, HandleMessage)
	b.Handle(tele.OnVoice, HandleMessage)
	b.Handle(tele.OnVideoNote, HandleMessage)
	b.Handle(tele.OnDocument, HandleMessage)
	b.Handle(tele.OnAnimation, HandleMessage)
	b.Handle(tele.OnLocation, HandleMessage)
	b.Handle(tele.OnSticker, HandleMessage)

	//special handler for albums
	b.Handle(tele.OnMedia, CheckForwardAlbum)

	b.Start()
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

func initializeSqlite() {
	var err error
	db, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		logToChat("Failed to connect database, replies will be unavailable" + "\r\n" + err.Error())
	}

	// Migrate the schema
	err = db.AutoMigrate(&InboxMessage{})
	if err != nil {
		logToChat("Migration failed, replies will be unavailable" + "\r\n" + err.Error())
	}
}
