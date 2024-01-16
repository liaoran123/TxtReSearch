package web

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"txtresearch/pubgo"
	"txtresearch/xbdb"
)

func Idxpfx(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("Idxpfx")

	params := getparas(req)
	const (
		tbname   = "c"
		idxfield = "s"
	)

	kw := params["kw"]
	top := params["top"]
	if top == "" {
		top = "11"
	}
	caid := params["caid"]
	ckft, _ := req.Cookie("ft")
	l := false
	if ckft != nil {
		l = ckft.Value == "1"
	}

	//自动转化参数值
	idxvalue := Table[string(tbname)].Ifo.FieldChByte(idxfield, kw)
	count := 0
	count, _ = strconv.Atoi(top)
	if count > 11 {
		count = 11
	}
	idx := NewidxExefunc(kw, caid, count, l)
	idx.r.WriteString("[")
	Table["c"].Select.WhereIdxLikeFun([]byte(idxfield), idxvalue, true, idx.add)
	jsonstr := idx.r.String()
	if idx.r != nil {
		idx.r.Reset()
		bufpool.Put(idx.r)
	}
	jsonstr = strings.Trim(jsonstr, ",") + "]"
	w.Write([]byte(jsonstr))
	//w.Write([]byte(strconv.Quote(ef.r.String()))) //必须使用strconv.Quote转义
}

type idxExefunc struct {
	kw    string
	caid  string
	count int //返回条数
	loop  int
	ift   bool
	r     *bytes.Buffer
}

func NewidxExefunc(kw, caid string, count int, ift bool) *idxExefunc {
	return &idxExefunc{
		kw:    kw,
		caid:  caid,
		count: count,
		ift:   ift,
		r:     bufpool.Get().(*bytes.Buffer),
	}
}
func (i *idxExefunc) add(k, v []byte) bool {
	//rd := xbdb.KVToRd(k, v, []int{})
	ks := xbdb.SplitRd(k) //bytes.Split(rd, []byte(xbdb.Split))

	//过滤非目录下
	if !strings.HasPrefix(string(ks[2]), i.caid) {
		return true
	}
	sectext := string(ks[1])
	sectext = pubgo.Fj(sectext, i.ift)
	if strings.Contains(sectext, " ") { //if strings.Contains(string(ks[0]), " ") {
		return true
	}
	//k1 := string(ks[1])
	qsectext := strconv.Quote(sectext)
	if !strings.Contains(qsectext, sectext) { //存在需要转义的，都过滤
		return true
	}
	if !strings.Contains(i.r.String(), qsectext) {
		i.r.WriteString("{\"key\":" + qsectext + "},")
	} else {
		return true
	}
	i.loop++
	return i.loop < i.count
}
