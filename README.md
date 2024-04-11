## Telegram bot for channel feedback

This bot will forward any incoming messages to selected chat (private, group or channel). 
Chats with topics probably won't work for now.

Features:
- forward all incoming messages to a specific chat, channel or dm
- forward media group as album instead of separate messages
- reply to any incoming messages (including albums)
- launch with `--confirmreceive` to let users know that their message is received

### Run using docker:

```shell
docker run -d --name tg-feedback-bot --pull=always \
  -e "BOT_TOKEN=bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11" \
  -e "CHAT_ID=123456" \
  ghcr.io/rebel028/tg-feedback-bot:latest
```

### Run with Go
**Prerequisites:** [Go 1.19](https://go.dev/dl/)

Clone this repository then

Linux:
```shell
cd tg-feedback-bot
go build -o tg-feedback-bot
tg-feedback-bot -token bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11 -chatid 123456
```

Windows:
```powershell
cd ./tg-feedback-bot
go build -o tg-feedback-bot.exe
./tg-feedback-bot.exe -token bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11 -chatid 123456
```

where `bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11` is your bot token and `123456` is chat id to forward messages to.

[How to obtain bot token?](https://core.telegram.org/bots/features#botfather)

[How to get chat id?](https://www.google.com/search?q=how+to+get+telegram+chat+id)
