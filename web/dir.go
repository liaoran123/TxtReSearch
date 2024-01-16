package web

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"txtresearch/gstr"
	"txtresearch/pubgo"
)

var TopDir []dirinfo

type dir struct {
	Pubpar           pubpar
	Id, Name, Titile string
	Dirs             []dirinfo
	Filelist         []fileinfo
}
type fileinfo struct {
	Fid, Name string
	Isdir     bool
}
type dirinfo struct {
	Did, Name string
}

// 打开txt文件
func Dir(w http.ResponseWriter, req *http.Request) {
	//id := req.URL.Query().Get("id")
	id := gstr.Do(req.URL.Path+"#", "/", "#", true, false)
	cfpath := Dirs[id]
	//fpath := filepath.Join(pubgo.GetCurrentAbPath(), "/skqs", cfpath)
	fpath := filepath.Join(cdir, "/skqs", cfpath)
	filelist, err := os.ReadDir(fpath)
	if err != nil {
		fmt.Println(err)
		//w.Write([]byte(err.Error()))
	}
	//客户端实现过于复杂，故在此实现

	rd := dir{}
	rd.Id = id
	lname := ""
	for _, v := range filelist {
		ff := fileinfo{}
		ff.Isdir = v.IsDir()
		//fmt.Printf("v.Name(): %v\n", v.Name())
		//lname = gstr.Do(v.Name()+"~", "\\", "~", true, false)
		lname = v.Name() //gstr.RStr(v.Name(), "-")
		ff.Fid = id + "-" + gstr.LStr(v.Name(), "-")
		ff.Fid = strings.Trim(ff.Fid, "-")
		if ff.Isdir {
			ff.Name = lname
		} else {
			ff.Name = gstr.LStr(lname, ".txt") //gstr.Do(lname, "-", ".txt", false, false)
		}
		ff.Name = strings.Replace(ff.Name, "龙藏", "龙藏(乾隆大藏经)", -1) //针对乾隆大藏经seo
		rd.Filelist = append(rd.Filelist, ff)
	}
	rd.Pubpar = newpubpar(req)
	rd.Dirs = Getdirs(cfpath)
	for _, v := range rd.Dirs {
		rd.Titile += v.Name + "."
	}
	rd.Titile = strings.Replace(rd.Titile, "龙藏", "龙藏(乾隆大藏经)", -1) //针对乾隆大藏经seo
	rd.Titile = "四库全书." + rd.Titile
	rd.Titile = strings.Trim(rd.Titile, ".")
	rd.Name = gstr.Do(rd.Titile+"-", ".", "-", true, false)
	pubgo.Tj.Brows("dir")
	//--组织模板数据
	tpl := cdir + "/tpl"
	if ConfigMap["tpl"] != nil {
		tpl = ConfigMap["tpl"].(string)
	}
	TemplatesFiles := []string{
		tpl + "/dir.html",
		tpl + "/pub/static.html",
		tpl + "/pub/header.html",
		tpl + "/pub/search.html",
		tpl + "/pub/gomove.html",
		tpl + "/pub/footer.html", // 多加的文件
	}

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf,
	}
	t, _ := template.New("dir.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--New("index.html") 的 index.html必须是TemplatesFiles第一个文件名

	err = t.ExecuteTemplate(w, "dir.html", rd)
	if err != nil {
		fmt.Println(req.URL.Path, err)
	}
}

// 一个相对路径分解为层递路径
func Getdirs(path string) (r []dirinfo) {
	ds := strings.Split(path, "\\")
	dis := ""
	for _, v := range ds {
		if strings.HasSuffix(strings.ToLower(v), ".txt") { //匹配txt文件
			continue
		}
		di := dirinfo{}
		dis += "-" + gstr.LStr(v, "-")
		di.Did = strings.Trim(dis, "-")
		di.Name = gstr.RStr(v, "-")
		r = append(r, di)
	}
	return r
}

// 打开顶层目录，以作搜索项目选择
func GetTopDir() {
	fpath := filepath.Join(cdir, "/skqs")
	filelist, err := os.ReadDir(fpath)
	if err != nil {
		fmt.Println(err)
		//w.Write([]byte(err.Error()))
	}
	for _, v := range filelist {
		if !v.IsDir() {
			continue
		}
		df := dirinfo{}
		df.Did = gstr.LStr(v.Name(), "-")
		df.Name = gstr.RStr(v.Name(), "-")
		TopDir = append(TopDir, df)
	}
	TopDir = append(TopDir, dirinfo{Did: "10000", Name: "淘宝"}) //内置淘宝搜索
}
