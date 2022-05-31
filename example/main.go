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
			msg := ctx.Value.(*tgba.Message)
			ctx.Caller.Send(tgba.NewMessage(msg.Chat.ID, args))
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
