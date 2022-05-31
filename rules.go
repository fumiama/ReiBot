package rei

import (
	"strings"

	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// PrefixRule check if the message has the prefix and trim the prefix
//
// 检查消息前缀
func PrefixRule(prefixes ...string) Rule {
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*tgba.Message)
		if !ok || msg.Text == "" { // 确保无空
			return false
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(msg.Text, prefix) {
				ctx.State["prefix"] = prefix
				arg := strings.TrimLeft(msg.Text[len(prefix):], " ")
				ctx.State["args"] = arg
				return true
			}
		}
		return false
	}
}
