package routers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"txtresearch/createidx"
	"txtresearch/gstr"
	"txtresearch/pubgo"
	"txtresearch/xbdb"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

func SearchApi(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	params := getparas(req)
	tbname := "c"
	params["kw"] = Sublen(params["kw"], 35) //最大长度35
	Se := NewSeExefunc(tbname, params["kw"], params["dir"], 21)
	asc := params["asc"] == "" //params["asc"]默认空值即true
	p := params["p"]
	ok := false

	var key []byte
	var iter iterator.Iterator
	if p == "" {
		//第一页搜索没有p值，需要查询获得
		key = Table[tbname].Select.GetIdxPrefixLike([]byte("s"), []byte(Se.mkw))
		iter, ok = Table[tbname].Select.IterPrefixMove(key, asc)
	} else {
		iter, ok = Table[tbname].Select.IterSeekMove([]byte(p))
	}
	if !ok {
		return
	}
	ts := pubgo.Newts() //计算执行时间
	Se.r.WriteString("{\"datalist\":[")
	xbdb.NewIters(iter, ok, asc, 0, -1).ForKVFun(Se.search)
	jsonstr := Se.r.String()
	jsonstr = strings.Trim(jsonstr, ",")
	Se.r.Reset()
	Se.r.WriteString(jsonstr)
	Se.r.WriteString("],")
	Se.r.WriteString("\"lastkey\":" + strconv.Quote(Se.lastkey) + ",")
	setime := ts.Gstrts()
	Se.r.WriteString("\"setime\":\"" + setime + "\",")
	//fmt.Printf("setime: %v\n", setime)
	Se.r.WriteString("\"count\":" + strconv.Itoa(Se.loop) + "}")
	w.Write(Se.r.Bytes())
	//w.Write([]byte(strconv.Quote(Se.r.String()))) //必须使用strconv.Quote转义
	iter.Release()
	Se.r.Reset()
	bufpool.Put(Se.r)
}

// 搜索执行类
type SeExefunc struct {
	tbname  string
	kw      string
	dir     string   //目录范围，可以是多个
	ks      []string //用空格来判断组合查询，分解出多个词
	mkw     string   //最长的关键词
	count   int      //返回条数
	mlen    int      //摘录最大长度
	loop    int
	r       *bytes.Buffer
	lastkey string //最后的key值
	//--变量---
	//keys              [][]byte
	artid         string //文章id
	secid, lsecid string
	bt            time.Time
}

func NewSeExefunc(tbname, kw, dir string, count int) *SeExefunc {
	ks := strings.Split(kw, " ") //用空格来判断组合查询
	//获取字数最长的词，通常字数最长的就是数据量最少的词。以该词作为组合查询的遍历定位词。
	mkw := getMaxLenKw(ks)
	//maxkeylen := int(ConfigMap["maxkeylen"].(float64))
	mkw = Sublen(mkw, 7) //搜索词最大长度是7
	return &SeExefunc{
		tbname: tbname,
		kw:     kw,
		dir:    dir,
		ks:     ks,
		mkw:    mkw,
		count:  count,
		mlen:   49,
		r:      bufpool.Get().(*bytes.Buffer),
		bt:     time.Now(),
	}
}
func (e *SeExefunc) search(k, v []byte) bool {
	if time.Since(e.bt).Seconds() > 300000000000 { //只要是控制组合查询超时时间
		e.loop = 21 //以便用户点击下一页，分散时间进行搜索，缓解性能问题。
		//fmt.Println("组合查询超时3秒。") //多次执行由于会加载内存，则可以完成。
		return false
	}

	if !strings.Contains(string(k), e.mkw) {
		return false //key不存在kw，即已经超过所需数据
	}
	//rd := xbdb.KVToRd(k, v, []int{})
	//解构rd，转为字符串lastkey。参照artpost.ArtSecToId组合规则
	keys := Table[e.tbname].Split(k) //bytes.Split(rd, []byte(xbdb.Split))
	if len(keys) != 3 {
		fmt.Println("Split", string(k), "小于3个数字元素")
		return true
	}
	e.lastkey = string(k)
	artid := string(keys[2])
	secid := strings.Split(artid, ",")[1]
	artid = strings.Split(artid, ",")[0]
	if e.artid == artid { //排除重复。同一段落包含多个相同kw时，出现重复情况。
		return true
	}
	e.artid = artid
	e.artid = strings.Replace(e.artid, "|", "-", -1) //格式转换
	e.secid = secid
	if !strings.HasPrefix(artid, e.dir) {
		return true //范围搜索不匹配
	}
	if !e.exsit() { //组合查询
		return true
	}

	filepath := createidx.Dirs[e.artid] + ".txt"
	e.r.WriteString("{\"filepath\":" + strconv.Quote(filepath) + ",")                                //写入目录路径
	e.r.WriteString("\"fileid\":" + strconv.Quote(e.artid) + ",")                                    //写入目录路径
	e.r.WriteString("\"filename\":" + strconv.Quote(gstr.Do(filepath, "-", ".", true, false)) + ",") //写入文章标题
	e.r.WriteString("\"secid\":" + strconv.Quote(e.secid) + ",")
	e.r.WriteString("\"filemeta\":" + strconv.Quote(e.getartmeta(keys[2])) + ",") //写入文章摘录信息
	e.r.WriteString("\"lsecid\":" + strconv.Quote(e.lsecid) + "},")
	e.loop++
	return e.loop < e.count
}

// 组合查询
func (e *SeExefunc) exsit() (find bool) {
	if len(e.ks) < 2 {
		find = true
		return
	}
	id := e.artid
	idxvalue := Table[e.tbname].Ifo.FieldChByte("id", id)
	btext := Table[e.tbname].Select.GetPKValue(idxvalue)

	secstr := string(btext)
	fc := 0
	for i := 0; i < len(e.ks); i++ { //for _, v := range e.ks { //如果在该段落内容里，所有的词组都存在，即是匹配。
		if strings.Contains(secstr, Sublen(e.ks[i], 7)) {
			fc++
		}
		//find = find && strings.Contains(secstr, Sublen(e.ks[i], 7)) //精准查询
	}
	if fc >= len(e.ks)/2+1 { //存在一半以上即当为匹配
		find = true
	}
	return
}

// 文章摘录
func (e *SeExefunc) getartmeta(id []byte) (r string) {
	key := Table[e.tbname].Select.GetPkKey(id) //Table[e.tbname].Ifo.FieldChByte("id", skid)
	iter, ok := Table[e.tbname].Select.IterSeekMove(key)
	if !ok {
		return
	}
	meta := string(Table[e.tbname].Split(iter.Value())[0])
	//meta := string(iter.Value())
	for len(meta) < e.mlen*3 { //每个中文3个字节
		ok = iter.Next()
		if !ok {
			break
		}
		meta += "<br>" + string(Table[e.tbname].Split(iter.Value())[0])
	}

	/*
		[{"id":2,"title":"金刚经","fid":1,"isleaf":"0"},
		{"id":3,"title":"六祖坛经","fid":1,"isleaf":"0"}]
	*/
	keys := Table[e.tbname].Split(iter.Key()) //bytes.Split(rd, []byte(xbdb.Split))
	artid := string(keys[1])
	e.lsecid = strings.Split(artid, ",")[1]

	r = meta
	iter.Release()
	return
}

/*
// 写超时警告日志 通用方法

func TimeoutWarning(tag, detailed string, start time.Time, timeLimit float64) {
	dis := time.Now().Sub(start).Seconds()
	if dis > timeLimit {
		log.Warning(log.CENTER_COMMON_WARNING, tag, " detailed:", detailed, "TimeoutWarning using", dis, "s")
		//pubstr := fmt.Sprintf("%s count %v, using %f seconds", tag, count, dis)
		//stats.Publish(tag, pubstr)
	}
}
*/
//找出最大长度的词
func getMaxLenKw(ks []string) (s string) {
	l := 0
	lv := 0
	for _, v := range ks {
		lv = len([]rune(v))
		if lv >= l {
			s = v
			l = lv
		}
	}
	return
}
