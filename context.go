package rei

type Ctx struct {
	Event
	State
	Caller *TelegramClient
	ma     *Matcher
}
