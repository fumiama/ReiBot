package rei

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/sirupsen/logrus"
)

type Manager ctrl.Manager[*Ctx]

var (
	enmap     = make(map[string]*Engine)
	priomap   = make(map[int]string)    // priomap is map[prio]service
	foldermap = make(map[string]string) // foldermap is map[folder]service
	prio      uint64
	m         = ctrl.NewManager[*Ctx]("data/control/plugins.db", 10*time.Second)
)

// Register 注册插件控制器
func Register(service string, o *ctrl.Options[*Ctx]) *Engine {
	prio := int(atomic.AddUint64(&prio, 10))
	e := newEngine()
	s, ok := priomap[prio]
	if ok {
		panic(fmt.Sprint("prio", prio, "is used by", s))
	}
	priomap[prio] = service
	logrus.Debugln("[control]插件", service, "已设置优先级", prio)
	e.UsePreHandler(newctrl(service, o))
	e.prio = prio
	e.service = service
	switch {
	case o.PublicDataFolder != "":
		if unicode.IsLower([]rune(o.PublicDataFolder)[0]) {
			panic("public data folder " + o.PublicDataFolder + " must start with an upper case letter")
		}
		e.datafolder = "data/" + o.PublicDataFolder + "/"
	case o.PrivateDataFolder != "":
		if unicode.IsUpper([]rune(o.PrivateDataFolder)[0]) {
			panic("private data folder " + o.PrivateDataFolder + " must start with an lower case letter")
		}
		e.datafolder = "data/" + o.PrivateDataFolder + "/"
	default:
		e.datafolder = "data/rbp/"
	}
	if e.datafolder != "data/rbp/" {
		s, ok := foldermap[e.datafolder]
		if ok {
			panic("folder " + e.datafolder + " has been required by service " + s)
		}
		foldermap[e.datafolder] = service
	}
	if file.IsNotExist(e.datafolder) {
		err := os.MkdirAll(e.datafolder, 0755)
		if err != nil {
			panic(err)
		}
	}
	logrus.Debugln("[control]插件", service, "已设置数据目录", e.datafolder)
	enmap[service] = e
	return e
}

// Delete 删除插件控制器, 不会删除数据
func Delete(service string) {
	engine, ok := enmap[service]
	if ok {
		engine.Delete()
		m.RLock()
		_, ok = m.M[service]
		m.RUnlock()
		if ok {
			m.Lock()
			delete(m.M, service)
			m.Unlock()
		}
	}
}
