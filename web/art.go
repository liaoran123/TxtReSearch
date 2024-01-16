package web

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"txtresearch/gstr"
	"txtresearch/pubgo"
)

type art struct {
	Pubpar pubpar
	Aid    string
	Titile string
	Text   string
	IsEof  bool
	Pg     int
	Dirs   []dirinfo
}

// 建立一个共享进程池
var artpool = sync.Pool{
	New: func() interface{} {
		return new(art)
	},
}

// 打开txt文件
func Art(w http.ResponseWriter, req *http.Request) {
	//id := req.URL.Query().Get("id")
	id := gstr.Do(req.URL.Path+"#", "/", "#", true, false)
	arttjs.Brows(id)
	cfpath := Dirs[id]
	fpath := filepath.Join(pubgo.GetCurrentAbPath(), "/skqs", cfpath)
	pubgo.Tj.Brows("art")

	rd := artpool.Get().(*art)
	rd.Pubpar = newpubpar(req)
	rd.Aid = id
	rd.Titile = gstr.Do(fpath, "\\", ".txt", true, false)
	//rd.Titile = gstr.RStr(rd.Titile, "-")
	p := req.URL.Query().Get("p")
	intp, err := strconv.Atoi(p)
	if err != nil {
		intp = 1
	}
	rd.Pg = intp + 1
	rd.Text, rd.IsEof = opfile(fpath, intp) //string(filetext)
	//fmt.Println(rd.IsEof)
	rd.Dirs = Getdirs(cfpath)
	//--组织模板数据
	tpl := cdir + "/tpl"
	if ConfigMap["tpl"] != nil {
		tpl = ConfigMap["tpl"].(string)
	}
	TemplatesFiles := []string{
		tpl + "/art.html",
		tpl + "/pub/static.html",
		tpl + "/pub/header.html",
		tpl + "/pub/search.html",
		tpl + "/pub/gomove.html",
		tpl + "/pub/footer.html", // 多加的文件
	}

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf,
	}
	t, _ := template.New("art.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--New("index.html") 的 index.html必须是TemplatesFiles第一个文件名

	err = t.ExecuteTemplate(w, "art.html", rd)
	if err != nil {
		fmt.Println(req.URL.Path, err)
	}
	rd = new(art)   //重置struct
	artpool.Put(rd) //释放
}

func opfile(ph string, p int) (string, bool) {
	ctline := 49 //每次最多读取210行
	b := (p - 1) * ctline
	file, err := os.Open(ph)
	if err != nil {
		//log.Printf("Cannot open text file: %s, err: [%v]", textfile, err)
		return "", true
	}
	defer file.Close()
	loop := 0
	count := 0
	rtext := ""
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return rtext, true
		}
		if loop >= b {
			rtext += string(line) + "\n"
			count++
		}
		loop++
		if count >= ctline {
			return rtext, false
		}
	}
	/*
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text() // or
				//line := scanner.Bytes()
				if loop >= b {
					rtext += line
					count++
				}
				loop++
				if count > ctline {
					return rtext, false
				}
			}

		if err := scanner.Err(); err != nil {
			// log.Printf("Cannot scanner text file: %s, err: [%v]", textfile, err)
			return "", true
		}
	*/
	//return "", true
}
