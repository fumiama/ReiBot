<div align="center">
  <a href="https://crypko.ai/crypko/GtWYDpVMx5GYm/">
  <img src=".github/Misaki.png" alt="看板娘" width = "256">
  </a><br>

  <h1>ReiBot</h1>
  Lightweight Telegram bot framework<br><br>

  <img src="http://cmoe.azurewebsites.net/cmoe?name=ReiBot&theme=r34" /><br>

</div>

## Instructions

> Note: This framework is built mainly for Chinese users thus may display hard-coded Chinese prompts during the interaction.

This framework is a simple wrapper for [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api), aiming to make the event processing easier.

## Quick Start
> Here is a plugin-based example, see more in the `example` folder

![plugin-based example](https://user-images.githubusercontent.com/41315874/171567343-f61eba4e-2bc9-49b3-af05-6446f0a73c54.png)

```go
package main

import (
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	rei.OnMessagePrefix("echo").SetBlock(true).SecondPriority().
		Handle(func(ctx *rei.Ctx) {
			args := ctx.State["args"].(string)
			if args == "" {
				return
			}
			ctx.SendPlainMessage(false, args)
		})
	rei.Run(rei.Bot{
		Token:  "",
		Buffer: 256,
		UpdateConfig: tgba.UpdateConfig{
			Offset:  0,
			Limit:   0,
			Timeout: 60,
		},
		Debug: true,
	})
}
```

## Event-Based

> If Handler in Bot is implemented, the plugin function will be disabled.

![event-based example](https://user-images.githubusercontent.com/41315874/171567349-5ff59cfa-cc3a-44a8-8158-6c76c8d433b7.png)

```go
package main

import (
	"strings"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
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
		Handler: &rei.Handler{
			OnMessage: func(updateid int, bot *rei.TelegramClient, msg *tgba.Message) {
				if len(msg.Text) <= len("测试") {
					return
				}
				if !strings.HasPrefix(msg.Text, "测试") {
					return
				}
				_, err := bot.Send(tgba.NewMessage(msg.Chat.ID, msg.Text[len("测试"):]))
				if err != nil {
					log.Errorln(err)
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
					log.Errorln(err)
				}
			},
		},
	})
}
```

## Thanks

- [ZeroBot](https://github.com/wdvxdr1123/ZeroBot)
