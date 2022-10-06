package rei

// 生成空引擎
func newEngine() *Engine {
	return &Engine{
		preHandler:  []Rule{},
		midHandler:  []Rule{},
		postHandler: []Process{},
	}
}

var defaultEngine = newEngine()

// Engine is the pre_handler, mid_handler, post_handler manager
type Engine struct {
	preHandler  []Rule
	midHandler  []Rule
	postHandler []Process
	matchers    []*Matcher
	prio        int
	service     string
	datafolder  string
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

// ApplySingle 应用反并发
func (e *Engine) ApplySingle(s *Single[int64]) *Engine {
	s.Apply(e)
	return e
}

// DataFolder 本插件数据目录, 默认 data/rbp/
func (e *Engine) DataFolder() string {
	return e.datafolder
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

// OnMessagePrefix 前缀触发器
func OnMessagePrefix(prefix string, rules ...Rule) *Matcher {
	return defaultEngine.OnMessagePrefix(prefix, rules...)
}

// OnMessagePrefix 前缀触发器
func (e *Engine) OnMessagePrefix(prefix string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{PrefixRule(prefix)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessageSuffix 后缀触发器
func OnMessageSuffix(suffix string, rules ...Rule) *Matcher {
	return defaultEngine.OnMessageSuffix(suffix, rules...)
}

// OnMessageSuffix 后缀触发器
func (e *Engine) OnMessageSuffix(suffix string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{SuffixRule(suffix)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessageCommand 命令触发器
func OnMessageCommand(commands string, rules ...Rule) *Matcher {
	return defaultEngine.OnMessageCommand(commands, rules...)
}

// OnMessageCommand 命令触发器
func (e *Engine) OnMessageCommand(commands string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{CommandRule(commands)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessageRegex 正则触发器
func OnMessageRegex(regexPattern string, rules ...Rule) *Matcher {
	return defaultEngine.OnMessageRegex(regexPattern, rules...)
}

// OnMessageRegex 正则触发器
func (e *Engine) OnMessageRegex(regexPattern string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{RegexRule(regexPattern)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessageKeyword 关键词触发器
func OnMessageKeyword(keyword string, rules ...Rule) *Matcher {
	return defaultEngine.OnMessageKeyword(keyword, rules...)
}

// OnMessageKeyword 关键词触发器
func (e *Engine) OnMessageKeyword(keyword string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{KeywordRule(keyword)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessageFullMatch 完全匹配触发器
func OnMessageFullMatch(src string, rules ...Rule) *Matcher {
	return defaultEngine.OnMessageFullMatch(src, rules...)
}

// OnMessageFullMatch 完全匹配触发器
func (e *Engine) OnMessageFullMatch(src string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{FullMatchRule(src)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessageFullMatchGroup 完全匹配触发器组
func OnMessageFullMatchGroup(src []string, rules ...Rule) *Matcher {
	return defaultEngine.OnMessageFullMatchGroup(src, rules...)
}

// OnMessageFullMatchGroup 完全匹配触发器组
func (e *Engine) OnMessageFullMatchGroup(src []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{FullMatchRule(src...)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessageKeywordGroup 关键词触发器组
func OnMessageKeywordGroup(keywords []string, rules ...Rule) *Matcher {
	return defaultEngine.OnMessageKeywordGroup(keywords, rules...)
}

// OnMessageKeywordGroup 关键词触发器组
func (e *Engine) OnMessageKeywordGroup(keywords []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{KeywordRule(keywords...)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessageCommandGroup 命令触发器组
func OnMessageCommandGroup(commands []string, rules ...Rule) *Matcher {
	return defaultEngine.OnMessageCommandGroup(commands, rules...)
}

// OnMessageCommandGroup 命令触发器组
func (e *Engine) OnMessageCommandGroup(commands []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{CommandRule(commands...)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessagePrefixGroup 前缀触发器组
func OnMessagePrefixGroup(prefix []string, rules ...Rule) *Matcher {
	return defaultEngine.OnMessagePrefixGroup(prefix, rules...)
}

// OnMessagePrefixGroup 前缀触发器组
func (e *Engine) OnMessagePrefixGroup(prefix []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{PrefixRule(prefix...)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessageSuffixGroup 后缀触发器组
func OnMessageSuffixGroup(suffix []string, rules ...Rule) *Matcher {
	return defaultEngine.OnMessageSuffixGroup(suffix, rules...)
}

// OnMessageSuffixGroup 后缀触发器组
func (e *Engine) OnMessageSuffixGroup(suffix []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{SuffixRule(suffix...)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessageShell shell命令触发器
func OnMessageShell(command string, model interface{}, rules ...Rule) *Matcher {
	return defaultEngine.OnMessageShell(command, model, rules...)
}

// OnMessageShell shell命令触发器
func (e *Engine) OnMessageShell(command string, model interface{}, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "Message",
		Rules:  append([]Rule{ShellRule(command, model)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnCallbackQueryPrefix 前缀触发器
func OnCallbackQueryPrefix(prefix string, rules ...Rule) *Matcher {
	return defaultEngine.OnCallbackQueryPrefix(prefix, rules...)
}

// OnCallbackQueryPrefix 前缀触发器
func (e *Engine) OnCallbackQueryPrefix(prefix string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "CallbackQuery",
		Rules:  append([]Rule{PrefixRule(prefix)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnCallbackQuerySuffix 后缀触发器
func OnCallbackQuerySuffix(suffix string, rules ...Rule) *Matcher {
	return defaultEngine.OnCallbackQuerySuffix(suffix, rules...)
}

// OnCallbackQuerySuffix 后缀触发器
func (e *Engine) OnCallbackQuerySuffix(suffix string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "CallbackQuery",
		Rules:  append([]Rule{SuffixRule(suffix)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnCallbackQueryRegex 正则触发器
func OnCallbackQueryRegex(regexPattern string, rules ...Rule) *Matcher {
	return defaultEngine.OnCallbackQueryRegex(regexPattern, rules...)
}

// OnCallbackQueryRegex 正则触发器
func (e *Engine) OnCallbackQueryRegex(regexPattern string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "CallbackQuery",
		Rules:  append([]Rule{RegexRule(regexPattern)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnCallbackQueryKeyword 关键词触发器
func OnCallbackQueryKeyword(keyword string, rules ...Rule) *Matcher {
	return defaultEngine.OnCallbackQueryKeyword(keyword, rules...)
}

// OnCallbackQueryKeyword 关键词触发器
func (e *Engine) OnCallbackQueryKeyword(keyword string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "CallbackQuery",
		Rules:  append([]Rule{KeywordRule(keyword)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnCallbackQueryFullMatch 完全匹配触发器
func OnCallbackQueryFullMatch(src string, rules ...Rule) *Matcher {
	return defaultEngine.OnCallbackQueryFullMatch(src, rules...)
}

// OnCallbackQueryFullMatch 完全匹配触发器
func (e *Engine) OnCallbackQueryFullMatch(src string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "CallbackQuery",
		Rules:  append([]Rule{FullMatchRule(src)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnCallbackQueryFullMatchGroup 完全匹配触发器组
func OnCallbackQueryFullMatchGroup(src []string, rules ...Rule) *Matcher {
	return defaultEngine.OnCallbackQueryFullMatchGroup(src, rules...)
}

// OnCallbackQueryFullMatchGroup 完全匹配触发器组
func (e *Engine) OnCallbackQueryFullMatchGroup(src []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "CallbackQuery",
		Rules:  append([]Rule{FullMatchRule(src...)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnCallbackQueryKeywordGroup 关键词触发器组
func OnCallbackQueryKeywordGroup(keywords []string, rules ...Rule) *Matcher {
	return defaultEngine.OnCallbackQueryKeywordGroup(keywords, rules...)
}

// OnCallbackQueryKeywordGroup 关键词触发器组
func (e *Engine) OnCallbackQueryKeywordGroup(keywords []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "CallbackQuery",
		Rules:  append([]Rule{KeywordRule(keywords...)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnCallbackQueryPrefixGroup 前缀触发器组
func (e *Engine) OnCallbackQueryPrefixGroup(prefix []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "CallbackQuery",
		Rules:  append([]Rule{PrefixRule(prefix...)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnCallbackQuerySuffixGroup 后缀触发器组
func OnCallbackQuerySuffixGroup(suffix []string, rules ...Rule) *Matcher {
	return defaultEngine.OnCallbackQuerySuffixGroup(suffix, rules...)
}

// OnCallbackQuerySuffixGroup 后缀触发器组
func (e *Engine) OnCallbackQuerySuffixGroup(suffix []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   "CallbackQuery",
		Rules:  append([]Rule{SuffixRule(suffix...)}, rules...),
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}
