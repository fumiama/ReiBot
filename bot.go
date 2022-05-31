// Package rei ReiBot created in 2022.5.31
package rei

import (
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot bot 的配置
type Bot struct {
	// Token bot 的 token
	//    see https://core.telegram.org/bots#3-how-do-i-create-a-bot
	Token string `json:"token"`
	// Buffer 控制消息队列的长度
	Buffer int `json:"buffer"`
	// UpdateConfig 配置消息获取
	tgba.UpdateConfig
	// Debug 控制调试信息的输出与否
	Debug bool `json:"debug"`
	// Handler 注册对各种事件的处理
	Handler Handler
	// handlers 方便调用的 handler
	handlers map[string]GeneralHandleType
}

// Start clients without blocking
func Start(bots ...Bot) {
	for _, c := range bots {
		tc := NewTelegramClient(&c)
		tc.Connect()
		go tc.Listen()
	}
}

// Run clients and block self in listening last one
func Run(bots ...Bot) {
	var tc TelegramClient
	switch len(bots) {
	case 0:
		return
	case 1:
		c := bots[0]
		tc = NewTelegramClient(&c)
	default:
		for _, c := range bots[:len(bots)-1] {
			tc := NewTelegramClient(&c)
			tc.Connect()
			go tc.Listen()
		}
		c := bots[len(bots)-1]
		tc = NewTelegramClient(&c)
	}
	tc.Connect()
	tc.Listen()
}
