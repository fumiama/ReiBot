package echo

import (
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	rei.OnMessagePrefix("echo").SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			args := ctx.State["args"].(string)
			if args == "" {
				return
			}
			msg := ctx.Value.(*tgba.Message)
			_, _ = ctx.Caller.Send(tgba.NewMessage(msg.Chat.ID, args))
		})
}
