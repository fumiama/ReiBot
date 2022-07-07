package rei

import (
	"reflect"
	"strings"

	base14 "github.com/fumiama/go-base16384"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// Event ...
type Event struct {
	// Type is the non-null field name in Update
	Type string
	// UpdateID is the update's unique identifier.
	UpdateID int
	// Value is the non-null field value in Update
	Value any
}

func (tc *TelegramClient) processEvent(update tgba.Update) {
	v := reflect.ValueOf(&update).Elem()
	t := reflect.ValueOf(&update).Elem().Type()
	for i := 1; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.IsZero() {
			continue
		}
		tp := t.Field(i).Name
		if tc.b.Handler == nil {
			matcherLock.RLock()
			n := len(matcherMap[tp])
			if n == 0 {
				matcherLock.RUnlock()
				continue
			}
			log.Debugln("pass", tp, "event to plugins")
			matchers := make([]*Matcher, n)
			copy(matchers, matcherMap[tp])
			matcherLock.RUnlock()
			ctx := &Ctx{
				Event: Event{
					Type:     tp,
					UpdateID: update.UpdateID,
					Value:    f.Interface(),
				},
				State:  State{},
				Caller: tc,
			}
			switch tp {
			case "Message":
				ctx.Message = (*tgba.Message)(f.UnsafePointer())
				log.Println("receive Message Text:", ctx.Message.Text)
			case "CallbackQuery":
				ctx.Message = (*tgba.CallbackQuery)(f.UnsafePointer()).Message
				log.Println("receive CallbackQuery Data:", (*tgba.CallbackQuery)(f.UnsafePointer()).Data)
			}
			go match(ctx, matchers)
			continue
		}
		h, ok := tc.b.handlers[tp]
		if !ok {
			continue
		}
		log.Println("process", tp, "event")
		go h(update.UpdateID, tc, f.UnsafePointer())
	}
}

func match(ctx *Ctx, matchers []*Matcher) {
	if ctx.Message != nil && ctx.Message.Text != "" { // 确保无空
		ctx.IsToMe = func(ctx *Ctx) bool {
			if ctx.Message.Chat.IsPrivate() {
				return true
			}
			name := ctx.Caller.Self.String()
			if strings.HasPrefix(ctx.Message.Text, name) {
				logrus.Debugln("[event] message before process:", ctx.Message.Text)
				if len(ctx.Message.Entities) > 0 {
					n := len(name)
					c := 0
					i := 0
					for _, e := range ctx.Message.Entities {
						c += e.Length
						if c >= n {
							break
						}
						i++
					}
					if i > 0 {
						switch {
						case c < n:
							ctx.Message.Entities = nil
						case c == n:
							if i+1 < len(ctx.Message.Entities) {
								ctx.Message.Entities = ctx.Message.Entities[i+1:]
								ctx.Message.Entities[0].Offset = 0
							}
						default:
							ctx.Message.Entities = ctx.Message.Entities[i:]
							ctx.Message.Entities[0].Length -= c - n
							ctx.Message.Entities[0].Offset = 0
						}
						if len(ctx.Message.Entities) > 1 {
							o := ctx.Message.Entities[0].Length
							for _, e := range ctx.Message.Entities[1:] {
								e.Offset = o
								o += e.Length
							}
						}
					}
				}
				ctx.Message.Text = ctx.Message.Text[len(name):]
				logrus.Debugln("[event] message after process:", ctx.Message.Text)
				return true
			}
			u16txt, err := base14.UTF82UTF16BE(helper.StringToBytes(ctx.Message.Text))
			if err != nil {
				return false
			}
			for i, e := range ctx.Message.Entities {
				if e.IsMention() && e.Length > 0 {
					a := 2 * (e.Offset + 1)
					b := 2 * (e.Offset + e.Length)
					if a < b && a < len(u16txt) && b <= len(u16txt) {
						n, err := base14.UTF16BE2UTF8(u16txt[a:b])
						if err != nil {
							continue
						}
						if helper.BytesToString(n) == name {
							logrus.Debugln("[event] message before process:", ctx.Message.Text)
							n, err = base14.UTF16BE2UTF8(append(u16txt[:2*e.Offset], u16txt[b:]...))
							if err != nil {
								continue
							}
							ctx.Message.Text = helper.BytesToString(n)
							o := e.Offset
							ctx.Message.Entities = append(ctx.Message.Entities[:i], ctx.Message.Entities[i+1:]...)
							for _, e1 := range ctx.Message.Entities[i:] {
								e1.Offset = o
								o += e1.Length
							}
							if ctx.Message.Text[0] == ' ' {
								n := 0
								for c := range ctx.Message.Text {
									if c == ' ' {
										n++
									} else {
										break
									}
								}
								ctx.Message.Text = ctx.Message.Text[n:]
								u16txt = u16txt[2*n:]
								for _, e1 := range ctx.Message.Entities {
									if e1.Offset >= n {
										e1.Offset -= n
									}
								}
							}
							if ctx.Message.Text[len(ctx.Message.Text)-1] == ' ' {
								n := 0
								for i := len(ctx.Message.Text) - 1; i >= 0; i-- {
									if ctx.Message.Text[i] == ' ' {
										n++
									} else {
										break
									}
								}
								ctx.Message.Text = ctx.Message.Text[:len(ctx.Message.Text)-n]
								if len(ctx.Message.Entities) > 0 {
									elast := ctx.Message.Entities[len(ctx.Message.Entities)-n]
									if elast.Offset+elast.Length == len(u16txt)/2 {
										if elast.Length > n {
											elast.Length -= n
										} else {
											ctx.Message.Entities = ctx.Message.Entities[:len(ctx.Message.Entities)-1]
										}
									}
								}
							}
							logrus.Debugln("[event] message after process:", ctx.Message.Text)
							return true
						}
					}
				}
			}
			return false
		}(ctx)
	}
loop:
	for _, matcher := range matchers {
		for k := range ctx.State { // Clear State
			delete(ctx.State, k)
		}
		matcherLock.RLock()
		m := matcher.copy()
		matcherLock.RUnlock()
		ctx.ma = m

		// pre handler
		if m.Engine != nil {
			for _, handler := range m.Engine.preHandler {
				if !handler(ctx) { // 有 pre handler 未满足
					continue loop
				}
			}
		}

		for _, rule := range m.Rules {
			if rule != nil && !rule(ctx) { // 有 Rule 的条件未满足
				continue loop
			}
		}

		// mid handler
		if m.Engine != nil {
			for _, handler := range m.Engine.midHandler {
				if !handler(ctx) { // 有 mid handler 未满足
					continue loop
				}
			}
		}

		if m.Process != nil {
			m.Process(ctx) // 处理事件
		}
		if matcher.Temp { // 临时 Matcher 删除
			matcher.Delete()
		}

		if m.Engine != nil {
			// post handler
			for _, handler := range m.Engine.postHandler {
				handler(ctx)
			}
		}

		if m.Block { // 阻断后续
			break loop
		}
	}
}
