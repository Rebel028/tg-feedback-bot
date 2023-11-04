package main

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sort"
	"time"
)

type MediaGroupData struct {
	createdTime time.Time
	fileIds     []string
}

func NewMediaGroupData(fileIdToAppend string) *MediaGroupData {
	return &MediaGroupData{
		createdTime: time.Now(),
		fileIds:     []string{fileIdToAppend},
	}
}

func HandleMediaGroup(cache *map[string]*MediaGroupData, message *botApi.Message) {
	groupId := message.MediaGroupID
	fileId := GetFileId(message)
	log.Printf("Received media group item from %s, group Id %s", message.Chat.ID, groupId)
	if val, ok := (*cache)[groupId]; ok {
		log.Printf("Group Id %s already has %d items", groupId, len(val.fileIds))
		val.fileIds = append(val.fileIds, fileId)
	} else {
		log.Printf("Creating new group...")
		(*cache)[groupId] = NewMediaGroupData(fileId)
	}
}

func GetFileId(message *botApi.Message) string {
	if message.Photo != nil {
		//todo
		sort.Slice(message.Photo, func(i, j int) bool {
			return message.Photo[i].Width > message.Photo[j].Width
		})
		return message.Photo[0].FileID
	}
	panic("Not implemented!")
}
