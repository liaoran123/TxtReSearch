package web

import (
	"html/template"
	"net/http"
	"os"
	"strings"
	"txtresearch/gstr"
	"txtresearch/pubgo"
)

type index struct {
	Pubpar pubpar
	Ahot   []arthot
	Sehot  []string
}

// 首页的热文列表
type arthot struct {
	Id, Title string
}

// 读取根目录下的txt文件
func Txt() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.URL.Path, ".txt") {
				f, _ := os.ReadFile(cdir + req.URL.Path)
				w.Write(f)
				return //这里退出，则断开DoMiddleware之后的链路。
			}
			//DoMiddleware实则将这个代码之前的代码叠加执行，之后再执行f http.HandlerFunc。
			f(w, req)
		}
	}
}

func WebIndex(w http.ResponseWriter, req *http.Request) {
	rd := index{}
	ctx := req.Context().Value("data").(map[string]interface{})
	rd.Pubpar = ctx["pubpar"].(pubpar) //读取上下文
	pubgo.Tj.Brows("Index")

	for k := range arttjs.tjm {
		ahot := arthot{}
		ahot.Id = k
		ahot.Title = gstr.Do(Dirs[k], "\\", ".txt", true, false)
		ahot.Title = gstr.RStr(ahot.Title, "-")
		rd.Ahot = append(rd.Ahot, ahot)
	}
	for k := range searchtjs.tjm {
		rd.Sehot = append(rd.Sehot, k)
	}
	//--组织模板数据
	tpl := cdir + "/tpl" //必须是绝对路径
	if ConfigMap["tpl"] != nil {
		tpl = ConfigMap["tpl"].(string)
	}
	TemplatesFiles := []string{
		tpl + "/index.html",
		tpl + "/pub/static.html",
		tpl + "/pub/header.html",
		tpl + "/pub/search.html",
		tpl + "/pub/gomove.html",
		tpl + "/pub/footer.html", // 多加的文件
	}

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf,
	}
	t, err := template.New("index.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--New("index.html") 的 index.html必须是TemplatesFiles第一个文件名
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	err = t.ExecuteTemplate(w, "index.html", rd)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}

func Index() http.HandlerFunc {
	return DoMiddleware(WebIndex, Txt(), Pubpars()) //用中间件解耦合为3个HandlerFunc
}

/*
func Index(w http.ResponseWriter, req *http.Request) {
	keep := req.Context().Value("do").(bool) //读取上下文内容
	if !keep {
		return
	}

	rd := Newindex(req)
	pubgo.Tj.Brows("Index")

	for k := range arttjs.tjm {
		ahot := arthot{}
		ahot.Id = k
		ahot.Title = gstr.Do(Dirs[k], "\\", ".txt", true, false)
		ahot.Title = gstr.RStr(ahot.Title, "-")
		rd.Ahot = append(rd.Ahot, ahot)
	}
	for k := range searchtjs.tjm {
		rd.Sehot = append(rd.Sehot, k)
	}
	//--组织模板数据
	tpl := cdir + "/tpl" //必须是绝对路径
	if ConfigMap["tpl"] != nil {
		tpl = ConfigMap["tpl"].(string)
	}
	TemplatesFiles := []string{
		tpl + "/index.html",
		tpl + "/pub/static.html",
		tpl + "/pub/header.html",
		tpl + "/pub/search.html",
		tpl + "/pub/gomove.html",
		tpl + "/pub/footer.html", // 多加的文件
	}

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf,
	}
	t, err := template.New("index.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--New("index.html") 的 index.html必须是TemplatesFiles第一个文件名
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	err = t.ExecuteTemplate(w, "index.html", rd)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}
*/
