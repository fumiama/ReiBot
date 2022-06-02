package rei

import tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Ctx struct {
	Event
	State
	Caller *TelegramClient
	ma     *Matcher
}

// CheckSession 判断会话连续性
func (ctx *Ctx) CheckSession() Rule {
	msg := ctx.Value.(*tgba.Message)
	return func(ctx2 *Ctx) bool {
		msg2, ok := ctx.Value.(*tgba.Message)
		if !ok || msg.From == nil || msg.Chat == nil || msg2.From == nil || msg2.Chat == nil { // 确保无空
			return false
		}
		return msg.From.ID == msg2.From.ID && msg.Chat.ID == msg2.Chat.ID
	}
}
