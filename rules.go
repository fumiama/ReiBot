package rei

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

// PrefixRule check if the text message has the prefix and trim the prefix
//
// 检查消息前缀
func PrefixRule(prefixes ...string) Rule {
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *tgba.Message:
			if msg.Text == "" { // 确保无空
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
		case *tgba.CallbackQuery:
			if msg.Data == "" {
				return false
			}
			for _, prefix := range prefixes {
				if strings.HasPrefix(msg.Data, prefix) {
					ctx.State["prefix"] = prefix
					arg := strings.TrimLeft(msg.Data[len(prefix):], " ")
					ctx.State["args"] = arg
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// SuffixRule check if the text message has the suffix and trim the suffix
//
// 检查消息后缀
func SuffixRule(suffixes ...string) Rule {
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *tgba.Message:
			if msg.Text == "" { // 确保无空
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
		case *tgba.CallbackQuery:
			if msg.Data == "" {
				return false
			}
			for _, suffix := range suffixes {
				if strings.HasSuffix(msg.Data, suffix) {
					ctx.State["suffix"] = suffix
					arg := strings.TrimRight(msg.Data[:len(msg.Data)-len(suffix)], " ")
					ctx.State["args"] = arg
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// CommandRule check if the message is a command and trim the command name
//
//	this rule only supports tgba.Message
func CommandRule(commands ...string) Rule {
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*tgba.Message)
		if !ok || msg.Text == "" { // 确保无空
			return false
		}
		msg.Text = strings.TrimSpace(msg.Text)
		if msg.Text == "" { // 确保无空
			return false
		}
		cmdMessage := ""
		args := ""
		switch {
		case ctx.IsToMe && msg.IsCommand():
			cmdMessage = msg.Command()
			args = msg.CommandArguments()
			logrus.Debugln("CommandRule: IsCommand:", cmdMessage, "args:", args)
		case strings.HasPrefix(msg.Text, "/"):
			cmdMessage, args, _ = strings.Cut(msg.Text, " ")
			cmdMessage, _, _ = strings.Cut(cmdMessage, "@")
			cmdMessage = cmdMessage[1:]
			logrus.Debugln("CommandRule: has/prefix:", cmdMessage, "args:", args)
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
		switch msg := ctx.Value.(type) {
		case *tgba.Message:
			if msg.Text == "" { // 确保无空
				logrus.Debugln("RegexRule: null message text")
				return false
			}
			if matched := regex.FindStringSubmatch(msg.Text); matched != nil {
				ctx.State["regex_matched"] = matched
				logrus.Debugln("RegexRule: match message text", matched)
				return true
			}
			logrus.Debugln("RegexRule: no match message")
			return false
		case *tgba.CallbackQuery:
			if msg.Data == "" {
				logrus.Debugln("RegexRule: null query data")
				return false
			}
			if matched := regex.FindStringSubmatch(msg.Data); matched != nil {
				ctx.State["regex_matched"] = matched
				logrus.Debugln("RegexRule: match query data", matched)
				return true
			}
			logrus.Debugln("RegexRule: no match query data")
			return false
		default:
			logrus.Debugln("RegexRule: stub type")
			return false
		}
	}
}

// ReplyRule check if the message is replying some message
//
//	this rule only supports tgba.Message
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
		switch msg := ctx.Value.(type) {
		case *tgba.Message:
			if msg.Text == "" { // 确保无空
				return false
			}
			for _, str := range src {
				if strings.Contains(msg.Text, str) {
					ctx.State["keyword"] = str
					return true
				}
			}
			return false
		case *tgba.CallbackQuery:
			if msg.Data == "" {
				return false
			}
			for _, str := range src {
				if strings.Contains(msg.Data, str) {
					ctx.State["keyword"] = str
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// FullMatchRule check if src has the same copy of the message
func FullMatchRule(src ...string) Rule {
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *tgba.Message:
			if msg.Text == "" { // 确保无空
				return false
			}
			for _, str := range src {
				if str == msg.Text {
					ctx.State["matched"] = msg.Text
					return true
				}
			}
			return false
		case *tgba.CallbackQuery:
			if msg.Data == "" {
				return false
			}
			for _, str := range src {
				if str == msg.Data {
					ctx.State["matched"] = msg.Data
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// ShellRule 定义shell-like规则
//
//	this rule only supports tgba.Message
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
//
//	this rule only supports tgba.Message
func OnlyToMe(ctx *Ctx) bool {
	return ctx.IsToMe
}

// CheckUser only triggered by specific person
func CheckUser(userID ...int64) Rule {
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *tgba.Message:
			if msg.From == nil { // 确保无空
				return false
			}
			for _, uid := range userID {
				if msg.From.ID == uid {
					return true
				}
			}
			return false
		case *tgba.CallbackQuery:
			if msg.From == nil {
				return false
			}
			for _, uid := range userID {
				if msg.From.ID == uid {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// CheckChat only triggered in specific chat
func CheckChat(chatID ...int64) Rule {
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *tgba.Message:
			if msg.Chat == nil { // 确保无空
				return false
			}
			for _, cid := range chatID {
				if msg.Chat.ID == cid {
					return true
				}
			}
			return false
		case *tgba.CallbackQuery:
			if msg.Message == nil || msg.Message.Chat == nil {
				return false
			}
			for _, cid := range chatID {
				if msg.Message.Chat.ID == cid {
					return true
				}
			}
			return false
		default:
			return false
		}
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
	switch msg := ctx.Value.(type) {
	case *tgba.Message:
		if msg.From == nil { // 确保无空
			return false
		}
		for _, su := range ctx.Caller.b.SuperUsers {
			if su == msg.From.ID {
				return true
			}
		}
		return false
	case *tgba.CallbackQuery:
		if msg.From == nil {
			return false
		}
		for _, su := range ctx.Caller.b.SuperUsers {
			if su == msg.From.ID {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// CreaterPermission only triggered by the group creater or higher permission
func CreaterPermission(ctx *Ctx) bool {
	if SuperUserPermission(ctx) {
		return true
	}
	switch msg := ctx.Value.(type) {
	case *tgba.Message:
		if msg.From == nil || msg.Chat == nil { // 确保无空
			return false
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
	case *tgba.CallbackQuery:
		if msg.From == nil || msg.Message == nil || msg.Message.Chat == nil {
			return false
		}
		m, err := ctx.Caller.GetChatMember(
			tgba.GetChatMemberConfig{
				ChatConfigWithUser: tgba.ChatConfigWithUser{
					ChatID: msg.Message.Chat.ID,
					UserID: msg.From.ID,
				},
			},
		)
		if err != nil {
			return false
		}
		return m.IsCreator()
	default:
		return false
	}
}

// AdminPermission only triggered by the group admins or higher permission
func AdminPermission(ctx *Ctx) bool {
	if SuperUserPermission(ctx) {
		return true
	}
	switch msg := ctx.Value.(type) {
	case *tgba.Message:
		if msg.From == nil || msg.Chat == nil { // 确保无空
			return false
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
	case *tgba.CallbackQuery:
		if msg.From == nil || msg.Message == nil || msg.Message.Chat == nil {
			return false
		}
		m, err := ctx.Caller.GetChatMember(
			tgba.GetChatMemberConfig{
				ChatConfigWithUser: tgba.ChatConfigWithUser{
					ChatID: msg.Message.Chat.ID,
					UserID: msg.From.ID,
				},
			},
		)
		if err != nil {
			return false
		}
		return m.IsCreator() || m.IsAdministrator()
	default:
		return false
	}
}

// UserOrGrpAdmin 允许用户单独使用或群管使用
func UserOrGrpAdmin(ctx *Ctx) bool {
	if OnlyPublic(ctx) {
		return AdminPermission(ctx)
	}
	return OnlyToMe(ctx)
}

// IsPhoto 消息是图片返回 true
//
//	this rule only supports tgba.Message
func IsPhoto(ctx *Ctx) bool {
	msg, ok := ctx.Value.(*tgba.Message)
	if !ok || len(msg.Photo) == 0 { // 确保无空
		return false
	}
	ctx.State["photos"] = msg.Photo
	return true
}

// MustProvidePhoto 消息不存在图片阻塞120秒至有图片，超时返回 false
//
//	this rule only supports tgba.Message
func MustProvidePhoto(needphohint, failhint string) Rule {
	return func(ctx *Ctx) bool {
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
		next := NewFutureEvent("Message", 999, false, ctx.CheckSession(), IsPhoto).Next()
		select {
		case <-time.After(time.Second * 120):
			if failhint != "" {
				_, _ = ctx.SendPlainMessage(true, failhint)
			}
			return false
		case newCtx := <-next:
			ctx.State["photos"] = newCtx.State["photos"]
			ctx.Event = newCtx.Event
			return true
		}
	}
}
