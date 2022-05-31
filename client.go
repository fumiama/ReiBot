package rei

import (
	"log"
	"time"

	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramClient ...
type TelegramClient struct {
	tgba.BotAPI
	b Bot
}

// NewTelegramClient ...
func NewTelegramClient(c *Bot) (tc TelegramClient) {
	tc.b = *c
	return
}

// Connect ...
func (tc *TelegramClient) Connect() {
	log.Println("[INFO] 开始尝试连接到Telegram服务器, token:", tc.b.Token)
	for {
		ba, err := tgba.NewBotAPI(tc.b.Token)
		if err != nil {
			log.Println("[WARN] 连接到Telegram服务器时出现错误:", err, ", token:", tc.b.Token)
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		tc.BotAPI = *ba
		tc.Debug = tc.b.Debug
		tc.Buffer = tc.b.Buffer
		break
	}
	log.Println("[INFO] 连接到Telegram服务器成功, token:", tc.b.Token)
}

// Listen 开始监听事件
func (tc *TelegramClient) Listen() {
	log.Println("[INFO] 开始监听", tc.Self.UserName, "的事件")
	for {
		updates := tc.GetUpdatesChan(tc.b.UpdateConfig)
		for update := range updates {
			tc.processEvent(update)
		}
		log.Println("[WARN] Telegram服务器连接断开...")
		time.Sleep(time.Millisecond * time.Duration(3))
		tc.Connect()
	}
}
