package rei

import (
	"log"
	"reflect"
	"unsafe"

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
	h := reflect.ValueOf(tc.b.Handler)
	for i := 1; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.IsZero() {
			continue
		}
		m := h.FieldByName("On" + t.Field(i).Name)
		if m.IsZero() {
			continue
		}
		log.Println("[INFO] processEvent call", "On"+t.Field(i).Name)
		handler := m.Interface()
		go (*(*GeneralHandleType)(unsafe.Add(unsafe.Pointer(&handler), unsafe.Sizeof(uintptr(0)))))(update.UpdateID, tc, f.UnsafePointer())
	}
}
