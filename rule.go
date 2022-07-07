package rei

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/FloatTech/zbputils/img/writer"
	"github.com/FloatTech/zbputils/process"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/wdvxdr1123/ZeroBot/extension"
)

func newctrl(service string, o *ctrl.Options[*Ctx]) Rule {
	c := m.NewControl(service, o)
	return func(ctx *Ctx) bool {
		ctx.State["manager"] = c
		var gid int64 = 0
		if !ctx.Message.Chat.IsPrivate() {
			gid = ctx.Message.Chat.ID
		}
		return c.Handler(uintptr(unsafe.Pointer(ctx)), gid, ctx.value.Elem().FieldByName("From").FieldByName("ID").Int())
	}
}

func Lookup(service string) (*ctrl.Control[*Ctx], bool) {
	return m.Lookup(service)
}

func init() {
	process.NewCustomOnce(&m).Do(func() {
		OnMessageCommandGroup([]string{
			"响应", "response", "沉默", "silence",
		}, UserOrGrpAdmin).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			grp := ctx.Message.Chat.ID
			if ctx.Message.Chat.IsPrivate() {
				// 个人用户
				grp = -ctx.Message.From.ID
			}
			msg := ""
			switch ctx.State["command"] {
			case "响应", "response":
				err := m.Response(grp)
				if err == nil {
					msg = ctx.Caller.Self.String() + "将开始在此工作啦~"
				} else {
					msg = "ERROR: " + err.Error()
				}
			case "沉默", "silence":
				err := m.Silence(grp)
				if err == nil {
					msg = ctx.Caller.Self.String() + "将开始休息啦~"
				} else {
					msg = "ERROR: " + err.Error()
				}
			default:
				msg = "ERROR: bad command\"" + fmt.Sprint(ctx.State["command"]) + "\""
			}
			_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, msg))
		})

		OnMessageCommandGroup([]string{
			"启用", "enable", "禁用", "disable",
		}, UserOrGrpAdmin).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "没有找到指定服务!"))
				return
			}
			grp := ctx.Message.Chat.ID
			if ctx.Message.Chat.IsPrivate() {
				// 个人用户
				grp = -ctx.Message.From.ID
			}
			if strings.Contains(model.Command, "启用") || strings.Contains(model.Command, "enable") {
				service.Enable(grp)
				if service.Options.OnEnable != nil {
					service.Options.OnEnable(ctx)
				} else {
					_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "已启用服务: "+model.Args))
				}
			} else {
				service.Disable(grp)
				if service.Options.OnDisable != nil {
					service.Options.OnDisable(ctx)
				} else {
					_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "已禁用服务: "+model.Args))
				}
			}
		})

		OnMessageCommandGroup([]string{
			"全局启用", "allenable", "全局禁用", "alldisable",
		}, OnlyToMe, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "没有找到指定服务!"))
				return
			}
			if strings.Contains(model.Command, "启用") || strings.Contains(model.Command, "enable") {
				service.Enable(0)
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "已全局启用服务: "+model.Args))
			} else {
				service.Disable(0)
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "已全局禁用服务: "+model.Args))
			}
		})

		OnMessageCommandGroup([]string{"还原", "reset"}, UserOrGrpAdmin).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "没有找到指定服务!"))
				return
			}
			grp := ctx.Message.Chat.ID
			if ctx.Message.Chat.IsPrivate() {
				// 个人用户
				grp = -ctx.Message.From.ID
			}
			service.Reset(grp)
			_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "已还原服务的默认启用状态: "+model.Args))
		})

		OnMessageCommandGroup([]string{
			"禁止", "ban", "允许", "permit",
		}, AdminPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 2 {
				service, ok := Lookup(args[0])
				if !ok {
					_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "没有找到指定服务!"))
					return
				}
				grp := ctx.Message.Chat.ID
				if ctx.Message.Chat.IsPrivate() {
					// 个人用户
					grp = -ctx.Message.From.ID
				}
				msg := "**" + args[0] + "报告**"
				issu := SuperUserPermission(ctx)
				if strings.Contains(model.Command, "允许") || strings.Contains(model.Command, "permit") {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if issu {
								service.Permit(uid, grp)
								msg += "\n+ 已允许" + usr
							} else {
								member, err := ctx.Caller.GetChatMember(tgba.GetChatMemberConfig{ChatConfigWithUser: tgba.ChatConfigWithUser{ChatID: ctx.Message.Chat.ID, UserID: uid}})
								if err == nil && !member.HasLeft() && !member.WasKicked() {
									service.Permit(uid, grp)
									msg += "\n+ 已允许" + usr
								} else {
									msg += "\nx " + usr + " 不在本群"
								}
							}
						}
					}
				} else {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if issu {
								service.Ban(uid, grp)
								msg += "\n- 已禁止" + usr
							} else {
								member, err := ctx.Caller.GetChatMember(tgba.GetChatMemberConfig{ChatConfigWithUser: tgba.ChatConfigWithUser{ChatID: ctx.Message.Chat.ID, UserID: uid}})
								if err == nil && !member.HasLeft() && !member.WasKicked() {
									service.Ban(uid, grp)
									msg += "\n- 已禁止" + usr
								} else {
									msg += "\nx " + usr + " 不在本群"
								}
							}
						}
					}
				}
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, msg))
				return
			}
			_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "参数错误!"))
		})

		OnMessageCommandGroup([]string{
			"全局禁止", "allban", "全局允许", "allpermit",
		}, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 2 {
				service, ok := Lookup(args[0])
				if !ok {
					_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "没有找到指定服务!"))
					return
				}
				msg := "**" + args[0] + "全局报告**"
				if strings.Contains(model.Command, "允许") || strings.Contains(model.Command, "permit") {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							service.Permit(uid, 0)
							msg += "\n+ 已允许" + usr
						}
					}
				} else {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							service.Ban(uid, 0)
							msg += "\n- 已禁止" + usr
						}
					}
				}
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, msg))
				return
			}
			_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "参数错误!"))
		})

		OnMessageCommandGroup([]string{
			"封禁", "block", "解封", "unblock",
		}, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 1 {
				msg := "**报告**"
				if strings.Contains(model.Command, "解") || strings.Contains(model.Command, "un") {
					for _, usr := range args {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if m.DoUnblock(uid) == nil {
								msg += "\n- 已解封" + usr
							}
						}
					}
				} else {
					for _, usr := range args {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if m.DoBlock(uid) == nil {
								msg += "\n+ 已封禁" + usr
							}
						}
					}
				}
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, msg))
				return
			}
			_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "参数错误!"))
		})

		OnMessageCommandGroup([]string{
			"改变默认启用状态", "allflip",
		}, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "没有找到指定服务!"))
				return
			}
			err := service.Flip()
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR: "+err.Error()))
				return
			}
			_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "已改变全局默认启用状态: "+model.Args))
		})

		OnMessageCommandGroup([]string{"用法", "usage"}, UserOrGrpAdmin).SetBlock(true).secondPriority().
			Handle(func(ctx *Ctx) {
				model := extension.CommandModel{}
				_ = ctx.Parse(&model)
				service, ok := Lookup(model.Args)
				if !ok {
					_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "没有找到指定服务!"))
					return
				}
				if service.Options.Help != "" {
					gid := ctx.Message.Chat.ID
					if ctx.Message.Chat.IsPrivate() {
						// 个人用户
						gid = -ctx.Message.From.ID
					}
					_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, service.EnableMarkIn(gid).String()+" "+service.String()))
				} else {
					_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "该服务无帮助!"))
				}
			})

		OnMessageCommandGroup([]string{"服务列表", "service_list"}, UserOrGrpAdmin).SetBlock(true).secondPriority().
			Handle(func(ctx *Ctx) {
				i := 0
				gid := ctx.Message.Chat.ID
				if ctx.Message.Chat.IsPrivate() {
					// 个人用户
					gid = -ctx.Message.From.ID
				}
				m.RLock()
				msg := make([]any, 1, len(m.M)*4+1)
				m.RUnlock()
				msg[0] = "--------服务列表--------\n发送\"/用法 name\"查看详情\n发送\"/响应\"启用会话"
				m.ForEach(func(key string, manager *ctrl.Control[*Ctx]) bool {
					i++
					msg = append(msg, "\n", i, ": ", manager.EnableMarkIn(gid), key)
					return true
				})
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, fmt.Sprint(msg...)))
			})

		OnMessageCommandGroup([]string{"服务详情", "service_detail"}, UserOrGrpAdmin).SetBlock(true).secondPriority().
			Handle(func(ctx *Ctx) {
				i := 0
				gid := ctx.Message.Chat.ID
				if ctx.Message.Chat.IsPrivate() {
					// 个人用户
					gid = -ctx.Message.From.ID
				}
				m.RLock()
				msgs := make([]any, 1, len(m.M)*7+1)
				m.RUnlock()
				msgs[0] = "---服务详情---\n"
				m.ForEach(func(key string, service *ctrl.Control[*Ctx]) bool {
					i++
					msgs = append(msgs, i, ": ", service.EnableMarkIn(gid), key, "\n", service, "\n\n")
					return true
				})
				img, err := text.Render(fmt.Sprint(msgs...), text.FontFile, 400, 20)
				if err != nil {
					logrus.Errorf("[control] %v", err)
				}
				data, cl := writer.ToBytes(img.Image())
				_, err = ctx.Caller.Send(tgba.NewPhoto(ctx.Message.Chat.ID, tgba.FileBytes{
					Name:  "服务详情",
					Bytes: data,
				}))
				cl()
				if err != nil {
					_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR: "+err.Error()))
					return
				}
			})
	})
}
