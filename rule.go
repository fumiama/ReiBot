package rei

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/process"
	ctrl "github.com/FloatTech/zbpctrl"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

func newctrl(service string, o *ctrl.Options[*Ctx]) Rule {
	c := m.NewControl(service, o)
	return func(ctx *Ctx) bool {
		ctx.State["manager"] = c
		return c.Handler(ctx.Message.Chat.ID, ctx.value.Elem().FieldByName("From").Elem().FieldByName("ID").Int())
	}
}

func Lookup(service string) (*ctrl.Control[*Ctx], bool) {
	return m.Lookup(service)
}

// respLimiterManager 请求响应限速器管理
//
//	每 1d 4次触发
var respLimiterManager = rate.NewManager[int64](time.Hour*24, 4)

func init() {
	process.NewCustomOnce(&m).Do(func() {
		OnMessageCommandGroup([]string{
			"响应", "response", "沉默", "silence",
		}, UserOrGrpAdmin).SetBlock(true).Limit(func(ctx *Ctx) *rate.Limiter {
			return respLimiterManager.Load(ctx.Message.Chat.ID)
		}).secondPriority().Handle(func(ctx *Ctx) {
			grp := ctx.Message.Chat.ID
			msg := ""
			switch ctx.State["command"] {
			case "响应", "response":
				if m.CanResponse(grp) {
					msg = ctx.Caller.Self.String() + "已经在工作了哦~"
					break
				}
				if SuperUserPermission(ctx) {
					err := m.Response(grp)
					if err == nil {
						msg = ctx.Caller.Self.String() + "将开始在此工作啦~"
					} else {
						msg = "ERROR: " + err.Error()
					}
					break
				}
				notify := &tgba.PhotoConfig{
					BaseFile: tgba.BaseFile{
						BaseChat: tgba.BaseChat{
							ReplyMarkup: tgba.NewInlineKeyboardMarkup(
								tgba.NewInlineKeyboardRow(
									tgba.NewInlineKeyboardButtonData(
										"同意",
										"respermit"+fmt.Sprintf("%016x", uint64(grp)),
									),
									tgba.NewInlineKeyboardButtonData(
										"拒绝",
										"resrefuse"+fmt.Sprintf("%016x", uint64(grp)),
									),
								),
							),
						},
						File: func() tgba.RequestFileData {
							if ctx.Message.Chat.Photo != nil {
								return tgba.FileID(ctx.Message.Chat.Photo.BigFileID)
							}
							p, err := ctx.Caller.GetUserProfilePhotos(tgba.NewUserProfilePhotos(ctx.Message.From.ID))
							if err == nil && len(p.Photos) > 0 {
								fp := p.Photos[0]
								return tgba.FileID(fp[len(fp)-1].FileID)
							}
							return nil
						}(),
					},
					Caption:   "主人, @" + ctx.Message.From.String() + " 请求响应~\n*ChatType*: " + ctx.Message.Chat.Type + "\n*ChatUserName*: " + ctx.Message.Chat.UserName + "\n*ChatID*: " + strconv.FormatInt(ctx.Message.Chat.ID, 10) + "\n*ChatTitle*: " + ctx.Message.Chat.Title + "\n*ChatDescription*: " + ctx.Message.Chat.Description,
					ParseMode: "Markdown",
				}
				for _, id := range ctx.Caller.b.SuperUsers {
					notify.ChatID = id
					_, _ = ctx.Caller.Send(notify)
				}
				msg = "已将响应请求发给主人了, 请耐心等待回应哦~"
			case "沉默", "silence":
				if !m.CanResponse(grp) {
					msg = ctx.Caller.Self.String() + "已经在休息了哦~"
					break
				}
				err := m.Silence(grp)
				if err == nil {
					msg = ctx.Caller.Self.String() + "将开始休息啦~"
				} else {
					msg = "ERROR: " + err.Error()
				}
				if SuperUserPermission(ctx) {
					break
				}
				notify := &tgba.PhotoConfig{
					BaseFile: tgba.BaseFile{
						File: func() tgba.RequestFileData {
							if ctx.Message.Chat.Photo != nil {
								return tgba.FileID(ctx.Message.Chat.Photo.BigFileID)
							}
							p, err := ctx.Caller.GetUserProfilePhotos(tgba.NewUserProfilePhotos(ctx.Message.From.ID))
							if err == nil && len(p.Photos) > 0 {
								fp := p.Photos[0]
								return tgba.FileID(fp[len(fp)-1].FileID)
							}
							return nil
						}(),
					},
					Caption:   "主人, @" + ctx.Message.From.String() + " 主动结束了响应~\n*ChatType*: " + ctx.Message.Chat.Type + "\n*ChatUserName*: " + ctx.Message.Chat.UserName + "\n*ChatID*: " + strconv.FormatInt(ctx.Message.Chat.ID, 10) + "\n*ChatTitle*: " + ctx.Message.Chat.Title + "\n*ChatDescription*: " + ctx.Message.Chat.Description,
					ParseMode: "Markdown",
				}
				for _, id := range ctx.Caller.b.SuperUsers {
					notify.ChatID = id
					_, _ = ctx.Caller.Send(notify)
				}
			default:
				msg = "ERROR: bad command\"" + fmt.Sprint(ctx.State["command"]) + "\""
			}
			_, _ = ctx.SendPlainMessage(false, msg)
		})

		OnCallbackQueryRegex(`^respermit([0-9a-f]{16})$`, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			grp, err := strconv.ParseUint(ctx.State["regex_matched"].([]string)[1], 16, 64)
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "ERROR: "+err.Error()))
				return
			}
			msg := ""
			err = m.Response(int64(grp))
			if err == nil {
				msg = ctx.Caller.Self.String() + "将开始在此工作啦~"
			} else {
				msg = "ERROR: " + err.Error()
			}
			_, err = ctx.Caller.Send(&tgba.MessageConfig{
				BaseChat: tgba.BaseChat{
					ChatID: int64(grp),
				},
				Text: msg,
			})
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "ERROR: "+err.Error()))
				return
			}
			_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "已发送"))
		})

		OnCallbackQueryRegex(`^resrefuse([0-9a-f]{16})$`, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			grp, err := strconv.ParseUint(ctx.State["regex_matched"].([]string)[1], 16, 64)
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "ERROR: "+err.Error()))
				return
			}
			_, err = ctx.Caller.Send(&tgba.MessageConfig{
				BaseChat: tgba.BaseChat{
					ChatID: int64(grp),
				},
				Text: "很遗憾, 因为各种原因, 您暂时未获使用权限呢",
			})
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "ERROR: "+err.Error()))
				return
			}
			_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "已发送"))
		})

		OnMessageCommandGroup([]string{
			"全局响应", "allresponse", "全局沉默", "allsilence",
		}, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			msg := ""
			cmd := ctx.State["command"].(string)
			switch {
			case strings.Contains(cmd, "响应") || strings.Contains(cmd, "response"):
				err := m.Response(0)
				if err == nil {
					msg = ctx.Caller.Self.String() + "将开始在此工作啦~"
				} else {
					msg = "ERROR: " + err.Error()
				}
			case strings.Contains(cmd, "沉默") || strings.Contains(cmd, "silence"):
				err := m.Silence(0)
				if err == nil {
					msg = ctx.Caller.Self.String() + "将开始休息啦~"
				} else {
					msg = "ERROR: " + err.Error()
				}
			default:
				msg = "ERROR: bad command\"" + cmd + "\""
			}
			_, _ = ctx.SendPlainMessage(false, msg)
		})

		OnMessageCommandGroup([]string{
			"启用", "enable", "禁用", "disable",
		}, UserOrGrpAdmin).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			grp := ctx.Message.Chat.ID
			if !m.CanResponse(grp) {
				return
			}
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
				return
			}
			if strings.Contains(model.Command, "启用") || strings.Contains(model.Command, "enable") {
				service.Enable(grp)
				if service.Options.OnEnable != nil {
					service.Options.OnEnable(ctx)
				} else {
					_, _ = ctx.SendPlainMessage(false, "已启用服务: ", model.Args)
				}
			} else {
				service.Disable(grp)
				if service.Options.OnDisable != nil {
					service.Options.OnDisable(ctx)
				} else {
					_, _ = ctx.SendPlainMessage(false, "已禁用服务: ", model.Args)
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
				_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
				return
			}
			if strings.Contains(model.Command, "启用") || strings.Contains(model.Command, "enable") {
				service.Enable(0)
				_, _ = ctx.SendPlainMessage(false, "已全局启用服务: ", model.Args)
			} else {
				service.Disable(0)
				_, _ = ctx.SendPlainMessage(false, "已全局禁用服务: ", model.Args)
			}
		})

		OnMessageCommandGroup([]string{"还原", "reset"}, UserOrGrpAdmin).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			grp := ctx.Message.Chat.ID
			if !m.CanResponse(grp) {
				return
			}
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
				return
			}
			service.Reset(grp)
			_, _ = ctx.SendPlainMessage(false, "已还原服务的默认启用状态: ", model.Args)
		})

		OnMessageCommandGroup([]string{
			"禁止", "ban", "允许", "permit",
		}, AdminPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			grp := ctx.Message.Chat.ID
			if !m.CanResponse(grp) {
				return
			}
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 2 {
				service, ok := Lookup(args[0])
				if !ok {
					_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
					return
				}
				grp := ctx.Message.Chat.ID
				msg := "*" + args[0] + "报告*"
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
				_, _ = ctx.SendPlainMessage(false, msg)
				return
			}
			_, _ = ctx.SendPlainMessage(false, "参数错误!")
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
					_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
					return
				}
				msg := "*" + args[0] + "全局报告*"
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
				_, _ = ctx.SendPlainMessage(false, msg)
				return
			}
			_, _ = ctx.SendPlainMessage(false, "参数错误!")
		})

		OnMessageCommandGroup([]string{
			"封禁", "block", "解封", "unblock",
		}, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 1 {
				msg := "*报告*"
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
				_, _ = ctx.SendPlainMessage(false, msg)
				return
			}
			_, _ = ctx.SendPlainMessage(false, "参数错误!")
		})

		OnMessageCommandGroup([]string{
			"改变默认启用状态", "allflip",
		}, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
				return
			}
			err := service.Flip()
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			_, _ = ctx.SendPlainMessage(false, "已改变全局默认启用状态: ", model.Args)
		})

		OnMessageCommandGroup([]string{"用法", "usage"}, UserOrGrpAdmin).SetBlock(true).secondPriority().
			Handle(func(ctx *Ctx) {
				model := extension.CommandModel{}
				_ = ctx.Parse(&model)
				service, ok := Lookup(model.Args)
				if !ok {
					_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
					return
				}
				if service.Options.Help != "" {
					gid := ctx.Message.Chat.ID
					_, _ = ctx.SendPlainMessage(false, service.EnableMarkIn(gid), " ", service)
				} else {
					_, _ = ctx.SendPlainMessage(false, "该服务无帮助!")
				}
			})

		OnMessageCommandGroup([]string{"服务列表", "service_list"}, UserOrGrpAdmin).SetBlock(true).secondPriority().
			Handle(func(ctx *Ctx) {
				gid := ctx.Message.Chat.ID
				m.RLock()
				msg := make([]any, 1, len(m.M)*4+1)
				m.RUnlock()
				msg[0] = "--------服务列表--------\n发送\"/用法 name\"查看详情\n发送\"/响应\"启用会话"
				ForEachByPrio(func(i int, service *ctrl.Control[*Ctx]) bool {
					msg = append(msg, "\n", i+1, ": ", service.EnableMarkIn(gid), service.Service)
					return true
				})
				_, _ = ctx.SendPlainMessage(false, msg...)
			})

		OnMessageCommandGroup([]string{"服务详情", "service_detail"}, UserOrGrpAdmin).SetBlock(true).secondPriority().
			Handle(func(ctx *Ctx) {
				gid := ctx.Message.Chat.ID
				m.RLock()
				msgs := make([]any, 1, len(m.M)*7+1)
				m.RUnlock()
				msgs[0] = "---服务详情---\n"
				ForEachByPrio(func(i int, service *ctrl.Control[*Ctx]) bool {
					msgs = append(msgs, i+1, ": ", service.EnableMarkIn(gid), service.Service, "\n", service, "\n\n")
					return true
				})
				_, _ = ctx.SendPlainMessage(false, msgs...)
			})
	})
}
