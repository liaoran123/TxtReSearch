package createidx

//根据文件夹dzj（乾隆大藏经文本文件夹）生成搜索索引等
import (
	"os"
	"path/filepath"
	"strings"
	"txtresearch/pubgo"
	"txtresearch/xbdb"
)

// 根据dzj生成目录map
var Dirs map[string]string
var Root, ph, name string

// 获取大藏经所有经文
func GetDirs() {
	Dirs = make(map[string]string)
	Root = pubgo.GetCurrentAbPath() + "\\dzj"
	filepath.Walk(Root, walkFunc)
}

// 建立全文考据级索引
func CreatIdx(tb map[string]*xbdb.Table) {
	df := NewFileText(tb)
	for k, v := range Dirs {
		df.CreateIdx(k, Root+v+".txt")
	}
}
func walkFunc(path string, info os.FileInfo, err error) error {
	if strings.HasSuffix(strings.ToLower(info.Name()), ".txt") { //匹配txt文件
		ph = strings.Replace(path, Root, "", -1)
		name = getdirid(ph)
		Dirs[name] = strings.Replace(ph, ".txt", "", -1)
	}
	return nil
}

// 提取所有文件夹的id组合成该文件id
// dzj\3-论\5-此土著述\151-缁门警训\090-诫观檀越四事从苦缘起出生法 id=3-5-151-090
func getdirid(path string) string {
	xg := strings.Split(path, "\\")
	var hg []string
	id := ""
	for _, xv := range xg {
		hg = strings.Split(xv, "-")
		id += hg[0] + "-"
	}
	return strings.Trim(id, "-")
}

/*
//查找空文件夹
func walkFunc1(path string, info os.FileInfo, err error) error {

	if info.IsDir() {
		files, _ := ioutil.ReadDir(path)
		if len(files) == 0 { //空文件夹
			fmt.Println(path)
			dloop++
		}
	}
	return nil
}//dloop=4650
*/
