package rei

import (
	"encoding/binary"
	"strings"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	binutils "github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/math"
	"github.com/FloatTech/zbputils/process"
	b14 "github.com/fumiama/go-base16384"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var startTime int64

func init() {
	// 插件冲突检测 会在本群发送一条消息并在约 1s 后撤回
	OnMessageFullMatch("插件冲突检测", OnlyGroup, AdminPermission, OnlyToMe).SetBlock(true).secondPriority().
		Handle(func(ctx *Ctx) {
			tok := genToken()
			if tok == "" || len([]rune(tok)) != 4 {
				return
			}
			startTime = time.Now().Unix()
			msg, err := ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "●cd"+tok))
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR: "+err.Error()))
				return
			}
			process.SleepAbout1sTo2s()
			_, _ = ctx.Caller.Send(tgba.NewDeleteMessage(ctx.Message.Chat.ID, msg.MessageID))
		})

	OnMessageRegex("^●cd([\u4e00-\u8e00]{4})$", OnlyGroup).SetBlock(true).secondPriority().
		Handle(func(ctx *Ctx) {
			if isValidToken(ctx.State["regex_matched"].([]string)[1]) {
				gid := ctx.Message.Chat.ID
				w := binutils.SelectWriter()
				m.ForEach(func(key string, manager *ctrl.Control[*Ctx]) bool {
					if manager.IsEnabledIn(gid) {
						w.WriteString("\xfe\xff")
						w.WriteString(key)
					}
					return true
				})
				if w.Len() > 2 {
					my, cl := binutils.OpenWriterF(func(wr *binutils.Writer) {
						wr.WriteString("●cd●")
						wr.WriteString(b14.EncodeString(w.String()[2:]))
					})
					binutils.PutWriter(w)
					msg, err := ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, binutils.BytesToString(my)))
					cl()
					if err != nil {
						_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR: "+err.Error()))
						return
					}
					process.SleepAbout1sTo2s()
					_, _ = ctx.Caller.Send(tgba.NewDeleteMessage(ctx.Message.Chat.ID, msg.MessageID))
				}
			}
		})

	OnMessageRegex("^●cd●(([\u4e00-\u8e00]*[\u3d01-\u3d06]?))", OnlyGroup).SetBlock(true).secondPriority().
		Handle(func(ctx *Ctx) {
			if time.Now().Unix()-startTime < 10 {
				gid := ctx.Message.Chat.ID
				for _, s := range strings.Split(b14.DecodeString(ctx.State["regex_matched"].([]string)[1]), "\xfe\xff") {
					m.RLock()
					c, ok := m.M[s]
					m.RUnlock()
					if ok && c.IsEnabledIn(gid) {
						c.Disable(gid)
					}
				}
			}
		})
}

func genToken() (tok string) {
	timebytes, cl := binutils.OpenWriterF(func(w *binutils.Writer) {
		w.WriteUInt64(uint64(time.Now().Unix()))
	})
	tok = b14.EncodeString(binutils.BytesToString(timebytes[1:]))
	cl()
	return
}

func isValidToken(tok string) (yes bool) {
	s := b14.DecodeString(tok)
	timebytes, cl := binutils.OpenWriterF(func(w *binutils.Writer) {
		_ = w.WriteByte(0)
		w.WriteString(s)
	})
	yes = math.Abs64(time.Now().Unix()-int64(binary.BigEndian.Uint64(timebytes))) < 10
	cl()
	return
}
