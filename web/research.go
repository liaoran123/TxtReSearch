package web

import (
	"html/template"
	"net/http"
	"txtresearch/pubgo"
)

type research struct {
	Pubpar pubpar
}

// 搜索说明
func WebResearch(w http.ResponseWriter, req *http.Request) {
	//数据组织
	rd := research{}
	ctx := req.Context().Value("data").(map[string]interface{})
	rd.Pubpar = ctx["pubpar"].(pubpar) //读取上下文

	pubgo.Tj.Brows("research")

	//--组织模板数据
	//--组织模板数据
	tpl := cdir + "/tpl"
	if ConfigMap["tpl"] != nil {
		tpl = ConfigMap["tpl"].(string)
	}
	TemplatesFiles := []string{
		tpl + "/research.html",
		tpl + "/pub/static.html",
		tpl + "/pub/header.html",
		tpl + "/pub/search.html",
		tpl + "/pub/gomove.html",
		tpl + "/pub/footer.html", // 多加的文件
	}

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf, //简繁转换
	}
	t, _ := template.New("research.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--  quanwen.html必须是TemplatesFiles第一个文件名
	t.ExecuteTemplate(w, "research.html", rd)

}
func Research() http.HandlerFunc {
	return DoMiddleware(WebResearch, Pubpars()) //用中间件解耦合为3个HandlerFunc
}
