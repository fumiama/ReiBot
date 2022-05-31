package rei

import (
	"unsafe"

	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GeneralHandleType func(int, *TelegramClient, unsafe.Pointer)

type Handler struct {
	OnMessage func(updateid int, bot *TelegramClient, msg *tgba.Message)

	OnEditedMessage func(updateid int, bot *TelegramClient, msg *tgba.Message)

	OnChannelPost func(updateid int, bot *TelegramClient, msg *tgba.Message)

	OnEditedChannelPost func(updateid int, bot *TelegramClient, msg *tgba.Message)

	OnInlineQuery func(updateid int, bot *TelegramClient, q *tgba.InlineQuery)

	OnChosenInlineResult func(updateid int, bot *TelegramClient, r *tgba.ChosenInlineResult)

	OnCallbackQuery func(updateid int, bot *TelegramClient, q *tgba.CallbackQuery)

	OnShippingQuery func(updateid int, bot *TelegramClient, q *tgba.ShippingQuery)

	OnPreCheckoutQuery func(updateid int, bot *TelegramClient, q *tgba.PreCheckoutQuery)

	OnPoll func(updateid int, bot *TelegramClient, p *tgba.Poll)

	OnPollAnswer func(updateid int, bot *TelegramClient, pa *tgba.PollAnswer)

	OnMyChatMember func(updateid int, bot *TelegramClient, m *tgba.ChatMemberUpdated)

	OnChatMember func(updateid int, bot *TelegramClient, m *tgba.ChatMemberUpdated)

	OnChatJoinRequest func(updateid int, bot *TelegramClient, r *tgba.ChatJoinRequest)
}
