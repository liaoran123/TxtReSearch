package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"txtresearch/pubgo"
	"txtresearch/xbdb"
)

type metas struct {
	Pubpar  pubpar
	Ahot    []arthot
	pgcount int
	Pg      []int
}

func (m *metas) getmetas(p int) {
	tbname := "meta"
	var b, e []byte
	bint := (p-1)*m.pgcount + 1
	eint := bint + m.pgcount
	b = Table[tbname].Ifo.FieldChByte("id", strconv.Itoa(bint))
	e = Table[tbname].Ifo.FieldChByte("id", strconv.Itoa(eint))
	b = Table[tbname].Select.GetPkKey(b)
	e = Table[tbname].Select.GetPkKey(e)
	Table[tbname].Select.FindRandFun(b, e, true, m.getkv)
}
func (m *metas) getkv(k, v []byte) bool {
	//println(string(k), string(v))
	vs := strings.Split(string(v), xbdb.Split)
	ahot := arthot{}
	ahot.Id = vs[1]
	ahot.Id = strings.Trim(ahot.Id, " ")
	ahot.Title = vs[2]
	if ahot.Title == "" {
		ahot.Title = "无题"
	}
	m.Ahot = append(m.Ahot, ahot)
	return true
}

// 打开txt文件
func Metas(w http.ResponseWriter, req *http.Request) {
	pubgo.Tj.Brows("metas")
	rd := new(metas)
	rd.Pubpar = newpubpar(req)
	rd.pgcount = 108 //每页49条记录
	p := req.URL.Query().Get("p")
	if p == "" {
		p = "1"
	}
	intp, err := strconv.Atoi(p)
	if err != nil {
		intp = 1
	}
	rd.getmetas(intp)               //string(filetext)
	if len(rd.Ahot) == rd.pgcount { //条数不够，则是没有了。
		for i := intp; i < intp+21; i++ {
			rd.Pg = append(rd.Pg, i)
		}
	}

	//--组织模板数据
	tpl := cdir + "/tpl"
	if ConfigMap["tpl"] != nil {
		tpl = ConfigMap["tpl"].(string)
	}
	TemplatesFiles := []string{
		tpl + "/metas.html",
		tpl + "/pub/static.html",
		tpl + "/pub/header.html",
		tpl + "/pub/search.html",
		tpl + "/pub/gomove.html",
		tpl + "/pub/footer.html", // 多加的文件
	}

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf,
	}
	t, _ := template.New("metas.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--New("index.html") 的 index.html必须是TemplatesFiles第一个文件名

	err = t.ExecuteTemplate(w, "metas.html", rd)
	if err != nil {
		fmt.Println(req.URL.Path, err)
	}
}
