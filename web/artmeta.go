package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
	"txtresearch/gstr"
	"txtresearch/pubgo"
)

type artmeta struct {
	Pubpar pubpar
	Aid    string
	Titile string
	Text   string
	IsEof  bool
	Pg     int
	Dirs   []dirinfo
}

// 建立一个共享进程池
var artmetapool = sync.Pool{
	New: func() interface{} {
		return new(artmeta)
	},
}

// 打开txt文件
func Artmeta(w http.ResponseWriter, req *http.Request) {
	id := gstr.Do(req.URL.Path+"#", "/", "#", true, false)
	kid := gstr.LStr(id, ",")
	cfpath := Dirs[kid]
	fpath := filepath.Join(pubgo.GetCurrentAbPath(), "/skqs", cfpath)
	pubgo.Tj.Brows("artmeta")

	rd := artmetapool.Get().(*artmeta)
	rd.Pubpar = newpubpar(req)
	rd.Aid = id
	rd.Titile = gstr.Do(fpath, "\\", ".txt", true, false)
	rd.Titile = gstr.RStr(rd.Titile, "-")

	rd.Text = getartmeta([]byte(id)) //string(filetext)
	//fmt.Println(rd.IsEof)
	rd.Dirs = Getdirs(cfpath)
	//--组织模板数据
	tpl := cdir + "/tpl"
	if ConfigMap["tpl"] != nil {
		tpl = ConfigMap["tpl"].(string)
	}
	TemplatesFiles := []string{
		tpl + "/artmeta.html",
		tpl + "/pub/static.html",
		tpl + "/pub/header.html",
		tpl + "/pub/search.html",
		tpl + "/pub/gomove.html",
		tpl + "/pub/footer.html", // 多加的文件
	}

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf,
	}
	t, _ := template.New("artmeta.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--New("index.html") 的 index.html必须是TemplatesFiles第一个文件名

	err := t.ExecuteTemplate(w, "artmeta.html", rd)
	if err != nil {
		fmt.Println(req.URL.Path, err)
	}
	rd = new(artmeta)   //重置struct
	artmetapool.Put(rd) //释放
}

func getartmeta(id []byte) (r string) {
	tbname := "c"
	key := Table[tbname].Select.GetPkKey(id) //Table[tbname].Ifo.FieldChByte("id", skid)
	iter, ok := Table[tbname].Select.IterSeekMove(key)
	if !ok {
		return
	}
	cv := string(Table[tbname].Split(iter.Value())[0])
	meta := ""
	hccount, loop := 0, 0
	for hccount >= 0 || loop > 108 { //连续2个回车处或超过108句
		meta += cv
		if cv == "\n" {
			hccount++
		}
		loop++
		ok = iter.Next()
		if !ok {
			break
		}
	}
	r = meta
	iter.Release()
	return
}
