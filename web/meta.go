package web

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"txtresearch/gstr"
	"txtresearch/pubgo"
	"txtresearch/xbdb"
)

type meta struct {
	Pubpar             pubpar
	Id                 string
	Titile             string
	Text               string
	text               strings.Builder //用这个方法拼接字符串速度快
	Aid, Atitle        string          //将文摘所在文章的id和标题
	Dirs               []dirinfo       //目录路径
	maxsec, hcct, loop int
}

func (m *meta) addtext(k, v []byte) bool {
	//println(string(k), string(v))
	m.loop++
	sec := strings.Split(string(v), xbdb.Split)[0]
	sec = strings.Trim(sec, " ")
	sec = strings.Trim(sec, "\u3000")
	if sec == "\n" {
		m.hcct++
	} else {
		m.hcct = 0
	}
	if m.Titile == "" {
		m.Titile = sec
		m.Titile = strings.Trim(m.Titile, "\n")
	}
	//m.Text += sec
	m.text.WriteString(sec)           //替换m.Text += sec，速度更快
	if m.loop >= 108 || m.hcct >= 2 { //最多108句
		return false
	}
	return true
}
func Meta(w http.ResponseWriter, req *http.Request) {
	const (
		tbname   = "c"
		idxfield = "id"
	)
	rd := new(meta)
	id := gstr.Do(req.URL.Path+"#", "/", "#", true, false)
	rd.Aid = gstr.LStr(id, ",")
	maxsec := req.URL.Query().Get("maxsec")
	if maxsec == "" {
		rd.maxsec = 108
	} else {
		rd.maxsec = 21
	}
	cfpath := Dirs[rd.Aid]
	fpath := filepath.Join(cdir, "/skqs", cfpath)
	rd.Dirs = Getdirs(cfpath)
	rd.Id = id
	rd.Pubpar = newpubpar(req)
	rd.Atitle = gstr.Do(fpath, "\\", ".txt", true, false)
	rd.Atitle = gstr.RStr(rd.Atitle, "-")

	pubgo.Tj.Brows("meta")

	key := Table[tbname].Ifo.FieldChByte(idxfield, id)
	key = Table[tbname].Select.GetPkKey(key)
	Table[tbname].Select.FindSeekFun(key, true, rd.addtext)
	rd.Text = rd.text.String()
	rd.text.Reset() //释放内存
	rd.Titile = strings.Trim(rd.Titile, "。")
	if len(rd.Titile) > 35*3 {
		rd.Titile = getlimitstr(rd.Titile)
	}

	//--组织模板数据
	tpl := cdir + "/tpl"
	if ConfigMap["tpl"] != nil {
		tpl = ConfigMap["tpl"].(string)
	}
	TemplatesFiles := []string{
		tpl + "/meta.html",
		tpl + "/pub/static.html",
		tpl + "/pub/header.html",
		tpl + "/pub/search.html",
		tpl + "/pub/gomove.html",
		tpl + "/pub/footer.html", // 多加的文件
	}

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf,
	}
	t, err := template.New("meta.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	err = t.ExecuteTemplate(w, "meta.html", rd)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func getlimitstr(text string) string {
	title := ""
	ss := pubgo.GetCnS(text)
	for _, sv := range ss {
		if strings.Trim(sv, " ") == "" {
			continue
		} else {
			title += sv + " "
			if len(title) >= 21 { //至少7个字
				break
			}
		}
	}
	return title
}
