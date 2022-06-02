package rei

import (
	"sort"
	"sync"

	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

type (
	// Rule filter the event
	Rule func(ctx *Ctx) bool
	// Process 事件处理函数
	Process func(ctx *Ctx)
)

// Matcher 是 ZeroBot 匹配和处理事件的最小单元
type Matcher struct {
	// Temp 是否为临时Matcher，临时 Matcher 匹配一次后就会删除当前 Matcher
	Temp bool
	// Block 是否阻断后续 Matcher，为 true 时当前Matcher匹配成功后，后续Matcher不参与匹配
	Block bool
	// priority 优先级，越小优先级越高
	priority int
	// Event 当前匹配到的事件
	Event *Event
	// Type 匹配的事件类型
	Type string
	// Rules 匹配规则
	Rules []Rule
	// Process 处理事件的函数
	Process Process
	// Engine 注册 Matcher 的 Engine，Engine可为一系列 Matcher 添加通用 Rule 和 其他钩子
	Engine *Engine
}

var (
	// 所有主匹配器列表
	matcherMap = make(map[string][]*Matcher, 0)
	// Matcher 修改读写锁
	matcherLock = sync.RWMutex{}
)

// State store the context of a matcher.
type State map[string]interface{}

func sortMatcher(typ string) {
	sort.Slice(matcherMap[typ], func(i, j int) bool { // 按优先级排序
		return matcherMap[typ][i].priority < matcherMap[typ][j].priority
	})
}

// SetBlock 设置是否阻断后面的 Matcher 触发
func (m *Matcher) SetBlock(block bool) *Matcher {
	m.Block = block
	return m
}

// Limit 限速器
//    postfn 当请求被拒绝时的操作
func (m *Matcher) Limit(limiterfn func(*Ctx) *rate.Limiter, postfn ...func(*Ctx)) *Matcher {
	m.Rules = append(m.Rules, func(ctx *Ctx) bool {
		if limiterfn(ctx).Acquire() {
			return true
		}
		if len(postfn) > 0 {
			for _, fn := range postfn {
				fn(ctx)
			}
		}
		return false
	})
	return m
}

// StoreMatcher store a matcher to matcher list.
func StoreMatcher(m *Matcher) *Matcher {
	matcherLock.Lock()
	defer matcherLock.Unlock()
	matcherMap[m.Type] = append(matcherMap[m.Type], m)
	sortMatcher(m.Type)
	return m
}

// StoreTempMatcher store a matcher only triggered once.
func StoreTempMatcher(m *Matcher) *Matcher {
	m.Temp = true
	StoreMatcher(m)
	return m
}

// Delete remove the matcher from list
func (m *Matcher) Delete() {
	matcherLock.Lock()
	defer matcherLock.Unlock()
	for i, matcher := range matcherMap[m.Type] {
		if m == matcher {
			matcherMap[m.Type] = append(matcherMap[m.Type][:i], matcherMap[m.Type][i+1:]...)
		}
	}
}

func (m *Matcher) copy() *Matcher {
	return &Matcher{
		Type:     m.Type,
		Rules:    m.Rules,
		Block:    m.Block,
		priority: m.priority,
		Process:  m.Process,
		Temp:     m.Temp,
		Engine:   m.Engine,
	}
}

// Handle 直接处理事件
func (m *Matcher) Handle(handler Process) *Matcher {
	m.Process = handler
	return m
}
