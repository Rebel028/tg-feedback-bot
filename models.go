package main

import (
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type ConfigParams struct {
	botToken          string
	forwardMessagesTo int64
	confirmReceive    bool
	useLogger         bool
}

type MediaGroupItem struct {
	updateId    int
	media       tele.Media
	messageId   int
	messageText string
}

type InboxMessage struct {
	ID           uint
	MessageId    int `gorm:"index"`
	OriginalId   int
	OriginalChat int64
}

type SuggestedPostMetadata struct {
	ID       string `gorm:"primaryKey"`
	SenderId int64
	Files    []FileMetadata
	Posted   gorm.DeletedAt
}

type FileMetadata struct {
	ID                      uint
	SuggestedPostMetadataId string
	FileId                  string
	Type                    string
}
