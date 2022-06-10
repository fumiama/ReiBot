package rei

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// PrefixRule check if the text message has the prefix and trim the prefix
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

// SuffixRule check if the text message has the suffix and trim the suffix
//
// 检查消息后缀
func SuffixRule(suffixes ...string) Rule {
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*tgba.Message)
		if !ok || msg.Text == "" { // 确保无空
			return false
		}
		for _, suffix := range suffixes {
			if strings.HasSuffix(msg.Text, suffix) {
				ctx.State["suffix"] = suffix
				arg := strings.TrimRight(msg.Text[:len(msg.Text)-len(suffix)], " ")
				ctx.State["args"] = arg
				return true
			}
		}
		return false
	}
}

// CommandRule check if the message is a command and trim the command name
func CommandRule(commands ...string) Rule {
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*tgba.Message)
		if !ok || msg.Text == "" { // 确保无空
			return false
		}
		cmdMessage := ""
		args := ""
		switch {
		case msg.IsCommand():
			cmdMessage = msg.Command()
			args = msg.CommandArguments()
		case strings.HasPrefix(msg.Text, "/"):
			a := strings.Index(msg.Text, "@")
			b := strings.Index(msg.Text, " ")
			switch {
			case b <= 1:
				cmdMessage = msg.Text[1:]
				args = ""
			case b == len(msg.Text):
				return false
			case a < 0:
				cmdMessage = msg.Text[1:b]
				args = msg.Text[b+1:]
			case a >= b:
				return false
			default:
				cmdMessage = msg.Text[1:a]
				args = msg.Text[b+1:]
			}
		default:
			return false
		}
		for _, command := range commands {
			if strings.HasPrefix(cmdMessage, command) {
				ctx.State["command"] = command
				ctx.State["args"] = args
				return true
			}
		}
		return false
	}
}

// RegexRule check if the message can be matched by the regex pattern
func RegexRule(regexPattern string) Rule {
	regex := regexp.MustCompile(regexPattern)
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*tgba.Message)
		if !ok || msg.Text == "" { // 确保无空
			return false
		}
		if matched := regex.FindStringSubmatch(msg.Text); matched != nil {
			ctx.State["regex_matched"] = matched
			return true
		}
		return false
	}
}

// ReplyRule check if the message is replying some message
func ReplyRule(messageID int) Rule {
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*tgba.Message)
		if !ok || msg.ReplyToMessage == nil { // 确保无空
			return false
		}
		return messageID == msg.MessageID
	}
}

// KeywordRule check if the message has a keyword or keywords
func KeywordRule(src ...string) Rule {
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*tgba.Message)
		if !ok || msg.Text == "" { // 确保无空
			return false
		}
		for _, str := range src {
			if strings.Contains(msg.Text, str) {
				ctx.State["keyword"] = str
				return true
			}
		}
		return false
	}
}

// FullMatchRule check if src has the same copy of the message
func FullMatchRule(src ...string) Rule {
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*tgba.Message)
		if !ok || msg.Text == "" { // 确保无空
			return false
		}
		for _, str := range src {
			if str == msg.Text {
				ctx.State["matched"] = msg.Text
				return true
			}
		}
		return false
	}
}

// ShellRule 定义shell-like规则
func ShellRule(cmd string, model interface{}) Rule {
	cmdRule := CommandRule(cmd)
	t := reflect.TypeOf(model)
	return func(ctx *Ctx) bool {
		if !cmdRule(ctx) {
			return false
		}
		// bind flag to struct
		args := ParseShell(ctx.State["args"].(string))
		val := reflect.New(t)
		fs := registerFlag(t, val)
		err := fs.Parse(args)
		if err != nil {
			return false
		}
		ctx.State["args"] = fs.Args()
		ctx.State["flag"] = val.Interface()
		return true
	}
}

// OnlyToMe only triggered in conditions of @bot or begin with the nicknames
func OnlyToMe(ctx *Ctx) bool {
	msg, ok := ctx.Value.(*tgba.Message)
	if !ok || msg.Text == "" { // 确保无空
		return false
	}
	if msg.Chat.IsPrivate() {
		return true
	}
	name := ctx.Caller.Self.String()
	if strings.HasPrefix(msg.Text, name) {
		return true
	}
	n := 0
	for _, e := range msg.Entities {
		if e.IsMention() && e.Length > 0 && msg.Text[n+1:n+e.Length] == name {
			return true
		}
		n += e.Length
	}
	return false
}

// CheckUser only triggered by specific person
func CheckUser(userId ...int64) Rule {
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*tgba.Message)
		if !ok || msg.From == nil { // 确保无空
			return false
		}
		for _, uid := range userId {
			if msg.From.ID == uid {
				return true
			}
		}
		return false
	}
}

// CheckChat only triggered in specific chat
func CheckChat(chatId ...int64) Rule {
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*tgba.Message)
		if !ok || msg.Chat == nil { // 确保无空
			return false
		}
		for _, cid := range chatId {
			if msg.Chat.ID == cid {
				return true
			}
		}
		return false
	}
}

// OnlyPrivate requires that the ctx.Event is private message
func OnlyPrivate(ctx *Ctx) bool {
	msg, ok := ctx.Value.(*tgba.Message)
	if !ok || msg.Chat == nil { // 确保无空
		return false
	}
	return msg.Chat.Type == "private"
}

// OnlyGroup requires that the ctx.Event is group message
func OnlyGroup(ctx *Ctx) bool {
	msg, ok := ctx.Value.(*tgba.Message)
	if !ok || msg.Chat == nil { // 确保无空
		return false
	}
	return msg.Chat.Type == "group"
}

// OnlySuperGroup requires that the ctx.Event is supergroup message
func OnlySuperGroup(ctx *Ctx) bool {
	msg, ok := ctx.Value.(*tgba.Message)
	if !ok || msg.Chat == nil { // 确保无空
		return false
	}
	return msg.Chat.Type == "supergroup"
}

// OnlyPublic requires that the ctx.Event is group or supergroup message
func OnlyPublic(ctx *Ctx) bool {
	msg, ok := ctx.Value.(*tgba.Message)
	if !ok || msg.Chat == nil { // 确保无空
		return false
	}
	return msg.Chat.Type == "supergroup" || msg.Chat.Type == "group"
}

// OnlyChannel requires that the ctx.Event is channel message
func OnlyChannel(ctx *Ctx) bool {
	msg, ok := ctx.Value.(*tgba.Message)
	if !ok || msg.Chat == nil { // 确保无空
		return false
	}
	return msg.Chat.Type == "channel"
}

// SuperUserPermission only triggered by the bot's owner
func SuperUserPermission(ctx *Ctx) bool {
	msg, ok := ctx.Value.(*tgba.Message)
	if !ok || msg.From == nil { // 确保无空
		return false
	}
	for _, su := range ctx.Caller.b.SuperUsers {
		if su == msg.From.ID {
			return true
		}
	}
	return false
}

// CreaterPermission only triggered by the group creater or higher permission
func CreaterPermission(ctx *Ctx) bool {
	msg, ok := ctx.Value.(*tgba.Message)
	if !ok || msg.From == nil || msg.Chat == nil { // 确保无空
		return false
	}
	for _, su := range ctx.Caller.b.SuperUsers {
		if su == msg.From.ID {
			return true
		}
	}
	m, err := ctx.Caller.GetChatMember(
		tgba.GetChatMemberConfig{
			ChatConfigWithUser: tgba.ChatConfigWithUser{
				ChatID: msg.Chat.ID,
				UserID: msg.From.ID,
			},
		},
	)
	if err != nil {
		return false
	}
	return m.IsCreator()
}

// AdminPermission only triggered by the group admins or higher permission
func AdminPermission(ctx *Ctx) bool {
	msg, ok := ctx.Value.(*tgba.Message)
	if !ok || msg.From == nil || msg.Chat == nil { // 确保无空
		return false
	}
	for _, su := range ctx.Caller.b.SuperUsers {
		if su == msg.From.ID {
			return true
		}
	}
	m, err := ctx.Caller.GetChatMember(
		tgba.GetChatMemberConfig{
			ChatConfigWithUser: tgba.ChatConfigWithUser{
				ChatID: msg.Chat.ID,
				UserID: msg.From.ID,
			},
		},
	)
	if err != nil {
		return false
	}
	return m.IsCreator() || m.IsAdministrator()
}

// UserOrGrpAdmin 允许用户单独使用或群管使用
func UserOrGrpAdmin(ctx *Ctx) bool {
	if OnlyPublic(ctx) {
		return AdminPermission(ctx)
	}
	return OnlyToMe(ctx)
}

// IsPhoto 消息是图片返回 true
func IsPhoto(ctx *Ctx) bool {
	msg, ok := ctx.Value.(*tgba.Message)
	if !ok || len(msg.Photo) == 0 { // 确保无空
		return false
	}
	ctx.State["photos"] = msg.Photo
	return true
}

// MustProvidePhoto 消息不存在图片阻塞120秒至有图片，超时返回 false
func MustProvidePhoto(ctx *Ctx, needphohint, failhint string) bool {
	msg, ok := ctx.Value.(*tgba.Message)
	if ok && len(msg.Photo) > 0 { // 确保无空
		ctx.State["photos"] = msg.Photo
		return true
	}
	// 没有图片就索取
	if needphohint != "" {
		_, err := ctx.Caller.Send(tgba.NewMessage(msg.Chat.ID, needphohint))
		if err != nil {
			return false
		}
	}
	next := NewFutureEvent("message", 999, false, ctx.CheckSession(), IsPhoto).Next()
	select {
	case <-time.After(time.Second * 120):
		if failhint != "" {
			_, _ = ctx.Caller.Send(tgba.NewMessage(msg.Chat.ID, failhint))
		}
		return false
	case newCtx := <-next:
		ctx.State["photos"] = newCtx.State["photos"]
		ctx.Event = newCtx.Event
		return true
	}
}
