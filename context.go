package rei

import (
	"fmt"
	"reflect"
	"sync"

	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Ctx struct {
	Event
	State
	Caller  *TelegramClient
	Message *tgba.Message
	ma      *Matcher
	IsToMe  bool
}

// decoder 反射获取的数据
type decoder []dec

type dec struct {
	index int
	key   string
}

// decoder 缓存
var decoderCache = sync.Map{}

// Parse 将 Ctx.State 映射到结构体
func (ctx *Ctx) Parse(model interface{}) (err error) {
	var (
		rv       = reflect.ValueOf(model).Elem()
		t        = rv.Type()
		modelDec decoder
	)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parse state error: %v", r)
		}
	}()
	d, ok := decoderCache.Load(t)
	if ok {
		modelDec = d.(decoder)
	} else {
		modelDec = decoder{}
		for i := 0; i < t.NumField(); i++ {
			t1 := t.Field(i)
			if key, ok := t1.Tag.Lookup("zero"); ok {
				modelDec = append(modelDec, dec{
					index: i,
					key:   key,
				})
			}
		}
		decoderCache.Store(t, modelDec)
	}
	for _, d := range modelDec { // decoder类型非小内存，无法被编译器优化为快速拷贝
		rv.Field(d.index).Set(reflect.ValueOf(ctx.State[d.key]))
	}
	return nil
}

// CheckSession 判断会话连续性
func (ctx *Ctx) CheckSession() Rule {
	msg := ctx.Value.(*tgba.Message)
	return func(ctx2 *Ctx) bool {
		msg2, ok := ctx.Value.(*tgba.Message)
		if !ok || msg.From == nil || msg.Chat == nil || msg2.From == nil || msg2.Chat == nil { // 确保无空
			return false
		}
		return msg.From.ID == msg2.From.ID && msg.Chat.ID == msg2.Chat.ID
	}
}

// Send 发送消息到对方
//
//	c.ChatID = ctx.Message.Chat.ID
//	if replytosender {
//		c.ReplyToMessageID = ctx.Message.MessageID
//	}
func (ctx *Ctx) Send(replytosender bool, c tgba.Chattable) (tgba.Message, error) {
	msg := reflect.ValueOf(c).Elem()
	msg.FieldByName("ChatID").SetInt(ctx.Message.Chat.ID)
	if replytosender {
		msg.FieldByName("ReplyToMessageID").SetInt(int64(ctx.Message.MessageID))
	}
	return ctx.Caller.Send(c)
}

// SendPlainMessage 发送无 entities 文本消息到对方
func (ctx *Ctx) SendPlainMessage(replytosender bool, printable ...any) (tgba.Message, error) {
	msg := tgba.NewMessage(ctx.Message.Chat.ID, fmt.Sprint(printable...))
	if replytosender {
		msg.ReplyToMessageID = ctx.Message.MessageID
	}
	return ctx.Caller.Send(&msg)
}

// SendMessage 发送富文本消息到对方
func (ctx *Ctx) SendMessage(replytosender bool, text string, entities ...tgba.MessageEntity) (tgba.Message, error) {
	msg := &tgba.MessageConfig{
		BaseChat: tgba.BaseChat{
			ChatID: ctx.Message.Chat.ID,
		},
		Text:     text,
		Entities: entities,
	}
	if replytosender {
		msg.ReplyToMessageID = ctx.Message.MessageID
	}
	return ctx.Caller.Send(msg)
}

// SendPhoto 发送图片消息到对方
func (ctx *Ctx) SendPhoto(file tgba.RequestFileData, replytosender bool, caption string, captionentities ...tgba.MessageEntity) (tgba.Message, error) {
	msg := &tgba.PhotoConfig{
		BaseFile: tgba.BaseFile{
			BaseChat: tgba.BaseChat{
				ChatID: ctx.Message.Chat.ID,
			},
			File: file,
		},
		Caption:         caption,
		CaptionEntities: captionentities,
	}
	if replytosender {
		msg.ReplyToMessageID = ctx.Message.MessageID
	}
	return ctx.Caller.Send(msg)
}

// SendAudio 发送音频消息到对方
func (ctx *Ctx) SendAudio(file tgba.RequestFileData, replytosender bool, caption string, captionentities ...tgba.MessageEntity) (tgba.Message, error) {
	msg := tgba.AudioConfig{
		BaseFile: tgba.BaseFile{
			BaseChat: tgba.BaseChat{
				ChatID: ctx.Message.Chat.ID,
			},
			File: file,
		},
		Caption:         caption,
		CaptionEntities: captionentities,
	}
	if replytosender {
		msg.ReplyToMessageID = ctx.Message.MessageID
	}
	return ctx.Caller.Send(&msg)
}
