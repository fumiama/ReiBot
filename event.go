package rei

import (
	"log"
	"reflect"

	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Event ...
type Event struct {
	// Type is the non-null field name in Update
	Type string
	// UpdateID is the update's unique identifier.
	UpdateID int
	// Value is the non-null field value in Update
	Value reflect.Value
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
		h, ok := tc.b.handlers[tp]
		if !ok {
			continue
		}
		log.Println("[INFO] process", tp, "event")
		go h(update.UpdateID, tc, f.UnsafePointer())
	}
}
