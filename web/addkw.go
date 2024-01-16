package web

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

var rw sync.RWMutex //读写锁

// 将内容追加文件末尾
func appendToFile(file, str string) {
	rw.Lock()
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		fmt.Printf("Cannot open file %s!\n", file)
		return
	}
	defer f.Close()
	f.WriteString("\n" + str)
	rw.Unlock()
}

// 判断文件是否大于某值，则写入表，后，清空
func writetb(file string, mlen int64) {
	rw.Lock()
	defer rw.Unlock()
	fi, err := os.Stat(file)
	if err != nil {
		return
	}
	if fi.Size() < mlen {
		return
	}
	f, err := os.Open(file)
	if err != nil {
		println(err.Error())
		return
	}
	defer f.Close()
	sline := ""
	reader := bufio.NewReader(f)
	insp := map[string]string{}

	loop, linelen := 0, 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		sline = string(line)
		linelen = len(sline)
		if linelen > 21 || linelen == 0 { //长度超过7，不添加
			continue
		}
		insp["id"] = sline
		r := Table["kw"].Ins(insp)
		if !r.Succ {
			println(r.Info)
		} else {
			loop++
		}
	}
	os.Truncate(file, 0) //清空文件
}
