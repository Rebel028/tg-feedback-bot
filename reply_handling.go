package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

var db *gorm.DB

func HandleMessage(c tele.Context) error {

	receivedMessage := c.Message()

	// if message is NOT from an admin(s) - forward it
	if receivedMessage.Chat.ID != forwardTo.ID {
		_, err := forwardMessage(b, forwardTo, receivedMessage)
		if err != nil {
			return err
		}
	} else {
		// if it's from admin(s), it's probably a reply if it's not we ignore it
		if !receivedMessage.IsReply() {
			return nil
		}

		// reply to message if we have messages db
		if db != nil {
			var replyFor InboxMessage
			if err := db.First(&replyFor, InboxMessage{MessageId: receivedMessage.ReplyTo.ID}).Error; err != nil {
				logToChat(err.Error())
				return err
			}
			chat, err := b.ChatByID(replyFor.OriginalChat)
			if err != nil {
				logToChat(err.Error())
				return err
			}
			replyTo := tele.Message{ID: replyFor.OriginalId, Chat: chat}
			sendable := makeSendable(receivedMessage)
			_, err = b.Reply(&replyTo, sendable)
			if err != nil {
				//if recipient deleted source message
				if strings.Contains(err.Error(), "message to reply not found") {
					_, err = b.Send(chat, sendable)
					if err == nil {
						return nil
					}
				}
				logToChat(err.Error())
				return err
			}
		}
	}

	return nil
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
	if _, err := forwardMessage(b, forwardTo, c.Message()); err != nil {
		logToChat(err.Error())
		return err
	}
	if confirmReceive {
		err := c.Reply("Message received!")
		if err != nil {
			logToChat(err.Error())
			return err
		}
	}
	return nil
}

func HandleAlbum(b *tele.Bot, c tele.Context, mediaGroupId string) {

	receivedAlbumPartMessage := c.Message()

	info := MediaGroupItem{
		media:       receivedAlbumPartMessage.Media(),
		updateId:    c.Update().ID,
		messageId:   receivedAlbumPartMessage.ID,
		messageText: receivedAlbumPartMessage.Caption,
	}

	log.Printf("Adding file to arr: %s", info.media.MediaFile().FileID)

	if messages, ok := mediaGroups[mediaGroupId]; ok {
		mediaGroups[mediaGroupId] = append(messages, info)
	} else {
		mediaGroups[mediaGroupId] = []MediaGroupItem{info}

		time.AfterFunc(time.Second, func() {
			album := tele.Album{}
			metadata := &SuggestedPostMetadata{
				SenderId: c.Sender().ID,
				Files:    make([]FileMetadata, 0),
			}

			// sort by update_id
			mediaGroup := mediaGroups[mediaGroupId]
			sort.Slice(mediaGroup, func(i, j int) bool {
				return mediaGroup[i].updateId < mediaGroup[j].updateId
			})

			caption := fmt.Sprintf("<b>Forwarded from</b> <a href=\"tg://user?id=%d\">%s</a>\n<b>Message</b>:\n%s",
				receivedAlbumPartMessage.Sender.ID,
				receivedAlbumPartMessage.Sender.FirstName+receivedAlbumPartMessage.Sender.LastName,
				mediaGroup[0].messageText)

			// collect media
			for i, albumInfo := range mediaGroup {
				media := albumInfo.media
				fileId := media.MediaFile().FileID
				log.Printf("Sending file %s", fileId)
				var c = ""
				if i == 0 {
					c = caption
				}
				album = append(album, createMedia(media, c))

				metadata.Files = append(metadata.Files, FileMetadata{
					FileId: fileId,
					Type:   media.MediaType(),
				})
			}

			results, err := b.SendAlbum(forwardTo, album, tele.ParseMode("Html"))
			if err != nil {
				logToChat(err.Error())
				return
			}

			if db != nil {
				// create a record for each received msg
				for i, albumMsg := range results {
					db.Create(&InboxMessage{
						MessageId:    albumMsg.ID,
						OriginalId:   mediaGroup[i].messageId,
						OriginalChat: receivedAlbumPartMessage.Chat.ID,
					})
				}
			}

			if confirmReceive {
				err := c.Reply("Message received!")
				if err != nil {
					logToChat(err.Error())
					return
				}
			}
		})
	}
}

func forwardMessage(b *tele.Bot, forwardTo *tele.Chat, receivedMessage *tele.Message) (*tele.Message, error) {
	forwardedMsg, err := b.Forward(forwardTo, receivedMessage)
	if err != nil {
		logToChat(err.Error())
		return nil, err
	}

	if db != nil {
		//create a record for received msg
		db.Create(&InboxMessage{
			MessageId:    forwardedMsg.ID,
			OriginalId:   receivedMessage.ID,
			OriginalChat: receivedMessage.Chat.ID,
		})
	}

	return forwardedMsg, err
}

func makeSendable(message *tele.Message) interface{} {
	if message.Media() == nil {
		return message.Text
	}
	return message.Media()
}

func HandleMessageEdited(context tele.Context) error {
	editedMessage := context.Message()

	// only handle edits from non-admin users
	if editedMessage.Chat.ID == forwardTo.ID {
		return nil
	}

	sender := editedMessage.Sender
	username := "@" + sender.Username
	if username == "" {
		username = sender.FirstName + " " + sender.LastName
	}

	text := fmt.Sprintf("%s UPD:\n%s", username, editedMessage.Text)
	_, err := b.Send(forwardTo, text)
	return err
}
