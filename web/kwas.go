package web

import (
	"html/template"
	"net/http"
	"txtresearch/gstr"
	"txtresearch/pubgo"
	"txtresearch/xbdb"
)

type kwas struct {
	Pubpar pubpar
	Kws    []string
	Chs    map[string]bool
}

func newkwas(req *http.Request) kwas {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return kwas{
		Pubpar: newpubpar(req),
		Chs:    map[string]bool{},
	}
}

func Kwas(w http.ResponseWriter, req *http.Request) {
	//数据组织
	rd := newkwas(req)
	pubgo.Tj.Brows("kwas")
	const (
		tbname   = "kw"
		idxfield = "id"
	)
	kw := gstr.Do(req.URL.Path+"#", "/", "#", true, false)
	idxvalue := Table[tbname].Ifo.FieldChByte(idxfield, kw)
	key := Table[tbname].Select.GetPkKey(idxvalue)
	r := Table[tbname].Select.FindPrefix(key, true, 0, 49, []int{}, false)
	//r := Table[tbname].Select.WhereIdxLike([]byte(idxfield), []byte(idxvalue), true, 0, 49, []int{}, false)
	if r != nil {
		for _, v := range r.Rd {
			rd.Kws = append(rd.Kws, string(xbdb.SplitRd(v)[0]))
		}
		for _, k := range rd.Kws {
			for _, c := range k {
				rd.Chs[string(c)] = true
			}
		}
	} else {
		for _, c := range kw {
			rd.Chs[string(c)] = true
		}
	}

	//--组织模板数据
	//--组织模板数据
	tpl := cdir + "/tpl"
	if ConfigMap["tpl"] != nil {
		tpl = ConfigMap["tpl"].(string)
	}
	TemplatesFiles := []string{
		tpl + "/kwas.html",
		tpl + "/pub/static.html",
		tpl + "/pub/header.html",
		tpl + "/pub/search.html",
		tpl + "/pub/gomove.html",
		tpl + "/pub/footer.html", // 多加的文件
	}

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf, //简繁转换
	}
	t, _ := template.New("kwas.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--  quanwen.html必须是TemplatesFiles第一个文件名
	t.ExecuteTemplate(w, "kwas.html", rd)

}
