package web

import (
	"os"
	"path/filepath"
	"strings"
	"txtresearch/gstr"
	"txtresearch/pubgo"
)

// 根据dzj生成目录map
var Dirs map[string]string
var Root, ph, name string
var dirtext string

// 获取大藏经所有经文
func GetDirs() {
	Root = pubgo.GetCurrentAbPath() + "\\skqs"
	//Root = pubgo.GetCurrentAbPath() + "\\skqs"
	filepath.Walk(Root, walkFunc)
	os.WriteFile("dirs.txt", []byte(dirtext), 0666)
}
func walkFunc(path string, info os.FileInfo, err error) error {
	ph = strings.Replace(path, Root, "", -1) //相对路径
	if ph == "" {
		return nil
	}
	ph = gstr.RStr(ph, "\\")
	name = pubgo.Getdirid(ph)
	dirtext += name + "," + ph + ";"
	return nil
}

/*
	func walkFunc1(path string, info os.FileInfo, err error) error {
		ph = strings.Replace(path, Root, "", -1) //相对路径
		if ph == "" {
			return nil
		}
		ph = gstr.RStr(ph, "\\")
		name = pubgo.Getdirid(ph)
		Dirs[name] = ph
		return nil
	}
*/
func LoadDirs() {
	text, err := os.ReadFile("dirs.txt")
	if err != nil {
		GetDirs()
	} else {
		dirtext = string(text)
	}
	dirs := strings.Split(dirtext, ";")
	Dirs = make(map[string]string)
	for _, v := range dirs {
		if v == "" {
			continue
		}
		dir := strings.Split(v, ",")
		Dirs[dir[0]] = dir[1]
	}
	dirtext = "" //清除内存
}

/*
函数参数说明 :

filename 操作的文件名

data 写入的内容

perm 文件不存在时创建文件并赋予的权限,例如 : 0666

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
*/
