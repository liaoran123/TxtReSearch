package web

import (
	"html/template"
	"net/http"
	"txtresearch/pubgo"
)

type shuoming struct {
	Pubpar pubpar
}

func newshuoming(req *http.Request) shuoming {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return shuoming{
		Pubpar: newpubpar(req),
	}
}

// 搜索说明
func Shuoming(w http.ResponseWriter, req *http.Request) {
	//数据组织
	rd := newshuoming(req)
	pubgo.Tj.Brows("shuoming")

	//--组织模板数据
	//--组织模板数据
	tpl := cdir + "/tpl"
	if ConfigMap["tpl"] != nil {
		tpl = ConfigMap["tpl"].(string)
	}
	TemplatesFiles := []string{
		tpl + "/shuoming.html",
		tpl + "/pub/static.html",
		tpl + "/pub/header.html",
		tpl + "/pub/search.html",
		tpl + "/pub/gomove.html",
		tpl + "/pub/footer.html", // 多加的文件
	}

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf, //简繁转换
	}
	t, _ := template.New("shuoming.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--  quanwen.html必须是TemplatesFiles第一个文件名
	t.ExecuteTemplate(w, "shuoming.html", rd)

}
