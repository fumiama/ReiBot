package main

import (
	_ "github.com/fumiama/ReiBot/example/echo"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	rei.OnMessageFullMatch("help").SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "echo string"))
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
