# ReiBot
Lightweight Telegram bot framework

## Instructions

This framework is a simple wrapper for [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api), aiming to make the event processing easier.

## Example

See under `example` folder or below.

![example](https://user-images.githubusercontent.com/41315874/171180885-c888a031-7797-4b4b-a232-9ff23f031b32.png)

```go
package main

import (
	"log"
	"strings"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	rei.Run(rei.Bot{
		Token:  "",
		Buffer: 256,
		UpdateConfig: tgba.UpdateConfig{
			Offset:  0,
			Limit:   0,
			Timeout: 60,
		},
		Debug: true,
		Handler: rei.Handler{
			OnMessage: func(updateid int, bot *rei.TelegramClient, msg *tgba.Message) {
				if len(msg.Text) <= len("测试") {
					return
				}
				if !strings.HasPrefix(msg.Text, "测试") {
					return
				}
				_, err := bot.Send(tgba.NewMessage(msg.Chat.ID, msg.Text[len("测试"):]))
				if err != nil {
					log.Println("[ERRO]", err)
				}
			},
			OnEditedMessage: func(updateid int, bot *rei.TelegramClient, msg *tgba.Message) {
				if len(msg.Text) <= len("测试") {
					return
				}
				if !strings.HasPrefix(msg.Text, "测试") {
					return
				}
				_, err := bot.Send(tgba.NewMessage(msg.Chat.ID, "已编辑："+msg.Text[len("测试"):]))
				if err != nil {
					log.Println("[ERRO]", err)
				}
			},
		},
	})
}
```