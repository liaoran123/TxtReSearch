package createidx

//根据文件夹dzj（乾隆大藏经文本文件夹）生成搜索索引等
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"txtresearch/gstr"
	"txtresearch/pubgo"
	"txtresearch/xbdb"
)

// 根据dzj生成目录map
var Dirs map[string]string
var Root, ph, name string
var txtfile []string

func getDirArr() {
	Dirs = make(map[string]string)
	Root = pubgo.GetCurrentAbPath() + "\\skqs"
	//Root = pubgo.GetCurrentAbPath() + "\\skqs"
	filepath.Walk(Root, walkFuncarr)
}

// 建立全文考据级索引
func CreatIdx(tb map[string]*xbdb.Table) {
	fmt.Printf("time.Now(): %v\n", time.Now())
	getDirArr()
	df := NewFileText(tb)
	for _, k := range txtfile {
		if v, ok := Dirs[k]; ok {
			df.CreateIdx(k, Root+"\\"+v+".txt")
		}
	}
	Dirs = make(map[string]string) //释放内存
	fmt.Printf("time.Now(): %v\n", time.Now())
}

// 所有txt文件
func walkFuncarr(path string, info os.FileInfo, err error) error {
	if strings.HasSuffix(strings.ToLower(info.Name()), ".txt") { //匹配txt文件
		ph = strings.Replace(path, Root, "", -1)
		ph = gstr.RStr(ph, "\\")
		name = pubgo.Getdirid(ph)
		Dirs[name] = strings.Replace(ph, ".txt", "", -1)
		txtfile = append(txtfile, name)
	}
	return nil
}

/*
// 提取所有文件夹的id组合成该文件id
// dzj\3-论\5-此土著述\151-缁门警训\090-诫观檀越四事从苦缘起出生法 id=3-5-151-090
func getdirid(path string) string {
	xg := strings.Split(path, "\\")
	var hg []string
	id := ""
	for _, xv := range xg {
		hg = strings.Split(xv, "-")
		if len(hg) < 2 {
			println(path)
		}
		id += hg[0] + "-"
	}
	return strings.Trim(id, "-")
}
*/
/*
// 移动文件

	func walkFunc1(path string, info os.FileInfo, err error) error {
		cp := path
		if info != nil {
			if info.IsDir() {
				files, _ := os.ReadDir(path)
				if len(files) != 0 { //非空文件夹
					for _, f := range files {
						//println(f.Name())
						if strings.HasSuffix(strings.ToLower(f.Name()), ".txt") { //匹配txt文件
							//移动文件
							sourcePath := filepath.Join(path, f.Name())
							destPath := filepath.Join(cp, f.Name())
							if sourcePath == destPath {
								break
							}
							err = os.Rename(sourcePath, destPath)
							if err != nil {
								fmt.Println(err)
							} else {
								fmt.Println(sourcePath)
								fmt.Println(destPath)
								fmt.Println("----------ok-----------------")
							}
						} else {
							//cp = path + f.Name()
							cp = filepath.Join(path, f.Name())
						}
					}
				}
			}
		}
		return nil
	} //dloop=4650
//改名
func walkFunc2(path string, info os.FileInfo, err error) error {
	if info != nil {
		if info.IsDir() {
			files, _ := os.ReadDir(path)
			fcount := len(files)
			ilen := len(strconv.Itoa(fcount))
			if fcount != 0 { //非空文件夹
				for i, f := range files {
					if strings.HasSuffix(strings.ToLower(f.Name()), ".txt") { //匹配txt文件
						sourcePath := filepath.Join(path, f.Name())
						destPath := filepath.Join(path, fmt.Sprintf("%0"+strconv.Itoa(ilen)+"d", i)+"-"+f.Name())
						err = os.Rename(sourcePath, destPath)
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println(sourcePath)
							fmt.Println(destPath)
							fmt.Println("----------ok-----------------")
						}
					} else {
						break
					}
				}
			}
		}
	}
	return nil
}
*/
