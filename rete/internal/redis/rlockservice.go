package redis

import (
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/rete/internal/types"
)

type lockServiceImpl struct {
	types.NwServiceImpl
	config common.Config
	key    string // key used to access lock
	redisutils.RedisHdl
	rand.Source
	sync.Mutex
	done chan bool
}

func NewLockServiceImpl(nw types.Network, config common.Config) types.LockService {
	r := lockServiceImpl{
		NwServiceImpl: types.NwServiceImpl{
			Nw: nw,
		},
		config:   config,
		RedisHdl: redisutils.NewRedisHdl(config.IDGens.Redis),
		Source:   rand.NewSource(time.Now().UnixNano()),
		done:     make(chan bool, 1),
	}
	return &r
}

func (l *lockServiceImpl) Init() {
	l.key = l.Nw.GetPrefix() + ":lock"
}

func (l *lockServiceImpl) Lock() {
	l.Mutex.Lock()
	value := strconv.FormatInt(l.Int63(), 10)
	for {
		ok, _ := l.Set(l.key, value, true, 16000)
		if ok {
			go func() {
				defer l.Mutex.Unlock()
				ticker := time.NewTicker(4 * time.Second)
				for {
					select {
					case <-ticker.C:
						l.Set(l.key, value, false, 16000)
					case <-l.done:
						ticker.Stop()
						l.DelIfEqual(l.key, value)
						return
					}
				}
			}()
			return
		}
		time.Sleep(128 * time.Millisecond)
	}
}

func (l *lockServiceImpl) Unlock() {
	l.done <- true
}
