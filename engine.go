package rei

// New 生成空引擎
func NewEngine() *Engine {
	return &Engine{
		preHandler:  []Rule{},
		midHandler:  []Rule{},
		postHandler: []Process{},
	}
}

var defaultEngine = NewEngine()

// Engine is the pre_handler, mid_handler, post_handler manager
type Engine struct {
	preHandler  []Rule
	midHandler  []Rule
	postHandler []Process
	matchers    []*Matcher
}

// Delete 移除该 Engine 注册的所有 Matchers
func (e *Engine) Delete() {
	for _, m := range e.matchers {
		m.Delete()
	}
}

// UsePreHandler 向该 Engine 添加新 PreHandler(Rule),
// 会在 Rule 判断前触发，如果 preHandler
// 没有通过，则 Rule, Matcher 不会触发
//
// 可用于分群组管理插件等
func (e *Engine) UsePreHandler(rules ...Rule) {
	e.preHandler = append(e.preHandler, rules...)
}

// UseMidHandler 向该 Engine 添加新 MidHandler(Rule),
// 会在 Rule 判断后， Matcher 触发前触发，如果 midHandler
// 没有通过，则 Matcher 不会触发
//
// 可用于速率限制等
func (e *Engine) UseMidHandler(rules ...Rule) {
	e.midHandler = append(e.midHandler, rules...)
}

// UsePostHandler 向该 Engine 添加新 PostHandler(Rule),
// 会在 Matcher 触发后触发，如果 PostHandler 返回 false,
// 则后续的 post handler 不会触发
//
// 可用于速率限制等
func (e *Engine) UsePostHandler(handler ...Process) {
	e.postHandler = append(e.postHandler, handler...)
}

// On 添加新的指定消息类型的匹配器(默认Engine)
func On(typ string, rules ...Rule) *Matcher { return defaultEngine.On(typ, rules...) }

// On 添加新的指定消息类型的匹配器
func (e *Engine) On(typ string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   typ,
		Rules:  rules,
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessage 消息触发器
func (e *Engine) OnMessage(rules ...Rule) *Matcher { return e.On("Message", rules...) }

// OnMessage 消息触发器
func OnMessage(rules ...Rule) *Matcher { return On("Message", rules...) }

// OnEditedMessage 修改消息触发器
func (e *Engine) OnEditedMessage(rules ...Rule) *Matcher { return e.On("EditedMessage", rules...) }

// OnEditedMessage 修改消息触发器
func OnEditedMessage(rules ...Rule) *Matcher { return On("EditedMessage", rules...) }

// OnChannelPost ...
func (e *Engine) OnChannelPost(rules ...Rule) *Matcher { return e.On("ChannelPost", rules...) }

// OnChannelPost ...
func OnChannelPost(rules ...Rule) *Matcher { return On("ChannelPost", rules...) }

// OnEditedChannelPost ...
func (e *Engine) OnEditedChannelPost(rules ...Rule) *Matcher {
	return e.On("EditedChannelPost", rules...)
}

// OnEditedChannelPost ...
func OnEditedChannelPost(rules ...Rule) *Matcher {
	return On("EditedChannelPost", rules...)
}

// OnInlineQuery ...
func (e *Engine) OnInlineQuery(rules ...Rule) *Matcher { return e.On("InlineQuery", rules...) }

// OnInlineQuery ...
func OnInlineQuery(rules ...Rule) *Matcher { return On("InlineQuery", rules...) }

// OnChosenInlineResult ...
func (e *Engine) OnChosenInlineResult(rules ...Rule) *Matcher {
	return e.On("ChosenInlineResult", rules...)
}

// OnChosenInlineResult ...
func OnChosenInlineResult(rules ...Rule) *Matcher { return On("ChosenInlineResult", rules...) }

// OnCallbackQuery ...
func (e *Engine) OnCallbackQuery(rules ...Rule) *Matcher { return e.On("CallbackQuery", rules...) }

// OnCallbackQuery ...
func OnCallbackQuery(rules ...Rule) *Matcher { return On("CallbackQuery", rules...) }

// OnShippingQuery ...
func (e *Engine) OnShippingQuery(rules ...Rule) *Matcher { return e.On("ShippingQuery", rules...) }

// OnShippingQuery ...
func OnShippingQuery(rules ...Rule) *Matcher { return On("ShippingQuery", rules...) }

// OnPreCheckoutQuery ...
func (e *Engine) OnPreCheckoutQuery(rules ...Rule) *Matcher {
	return e.On("PreCheckoutQuery", rules...)
}

// OnPreCheckoutQuery ...
func OnPreCheckoutQuery(rules ...Rule) *Matcher { return On("PreCheckoutQuery", rules...) }

// OnPoll ...
func (e *Engine) OnPoll(rules ...Rule) *Matcher { return e.On("Poll", rules...) }

// OnPoll ...
func OnPoll(rules ...Rule) *Matcher { return On("Poll", rules...) }

// OnPollAnswer ...
func (e *Engine) OnPollAnswer(rules ...Rule) *Matcher { return e.On("PollAnswer", rules...) }

// OnPollAnswer ...
func OnPollAnswer(rules ...Rule) *Matcher { return On("PollAnswer", rules...) }

// OnMyChatMember ...
func (e *Engine) OnMyChatMember(rules ...Rule) *Matcher { return e.On("MyChatMember", rules...) }

// OnMyChatMember ...
func OnMyChatMember(rules ...Rule) *Matcher { return On("MyChatMember", rules...) }

// OnChatMember ...
func (e *Engine) OnChatMember(rules ...Rule) *Matcher { return e.On("ChatMember", rules...) }

// OnChatMember ...
func OnChatMember(rules ...Rule) *Matcher { return On("ChatMember", rules...) }

// OnChatJoinRequest ...
func (e *Engine) OnChatJoinRequest(rules ...Rule) *Matcher { return e.On("ChatJoinRequest", rules...) }

// OnChatJoinRequest ...
func OnChatJoinRequest(rules ...Rule) *Matcher { return On("ChatJoinRequest", rules...) }

// OnPrefix 前缀触发器
func OnMessagePrefix(prefix string, rules ...Rule) *Matcher {
	return defaultEngine.OnMessagePrefix(prefix, rules...)
}

// OnPrefix 前缀触发器
func (e *Engine) OnMessagePrefix(prefix string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{PrefixRule(prefix)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}
