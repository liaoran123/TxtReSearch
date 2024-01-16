package web

import (
	"sync"
)

// ---统计文章，搜索流量
var cips, arttjs, searchtjs *tjs

type tjs struct {
	mlen int //最大map长度
	tjm  map[string]int
	mu   sync.RWMutex
}

func Newtjs(mlen int) *tjs {
	return &tjs{
		mlen: mlen,
		tjm:  make(map[string]int, mlen),
	}
}
func (t *tjs) Brows(id string) {
	if id == "" {
		return
	}
	if db, ok := t.tjm[id]; ok {
		t.mu.Lock()
		defer t.mu.Unlock()
		t.tjm[id] = db + 1
	} else {
		if len(t.tjm) >= t.mlen {
			t.Del()
		}
		t.tjm[id] = 1
	}
}

// 查找最小值
func (t *tjs) MinKey() string {
	mk := ""
	minct := 0
	for k, v := range t.tjm {
		if v == 1 { //1就是最小值
			return k
		}
		if minct == 0 { //初始化minct
			minct = v
		}
		if v <= minct {
			minct = v
			mk = k
		}
	}
	return mk
}
func (t *tjs) Del() {
	t.mu.Lock()
	defer t.mu.Unlock()
	k := t.MinKey()
	delete(t.tjm, k)
}
