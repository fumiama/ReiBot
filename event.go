package rei

import (
	"reflect"

	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
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
			log.Println("pass", tp, "event to plugins")
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
			case "CallbackQuery":
				ctx.Message = (*tgba.CallbackQuery)(f.UnsafePointer()).Message
			}
			match(ctx, matchers)
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
