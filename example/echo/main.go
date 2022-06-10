package echo

import (
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	rei.Register("echo", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help:             "- echo xxx",
	}).OnMessagePrefix("echo").SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			args := ctx.State["args"].(string)
			if args == "" {
				return
			}
			_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, args))
		})
}
