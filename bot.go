// Package rei ReiBot created in 2022.5.31
package rei

import (
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot bot 的配置
type Bot struct {
	// Token bot 的 token
	//    see https://core.telegram.org/bots#3-how-do-i-create-a-bot
	Token string
	// Buffer 控制消息队列的长度
	Buffer int
	// UpdateConfig 配置消息获取
	tgba.UpdateConfig
	// SuperUsers 超级用户
	SuperUsers []int64
	// Debug 控制调试信息的输出与否
	Debug bool
	// Handler 注册对各种事件的处理
	Handler *Handler
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

// GetBot 获取指定的bot (Ctx) 实例
func GetBot(id int64) *Ctx {
	caller, ok := clients.Load(id)
	if !ok {
		return nil
	}
	return &Ctx{Caller: caller}
}

// RangeBot 遍历所有bot (Ctx)实例
//
// 单次操作返回 true 则继续遍历，否则退出
func RangeBot(iter func(id int64, ctx *Ctx) bool) {
	clients.Range(func(key int64, value *TelegramClient) bool {
		return iter(key, &Ctx{Caller: value})
	})
}
