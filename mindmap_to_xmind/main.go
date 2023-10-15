package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/jan-bar/xmind"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := flag.String("db", "", "db path")
	file := flag.String("file", "", "mindmap file path")
	title := flag.String("title", "", "mindmap title")
	dst := flag.String("dst", "save.xmind", "save xmind path")
	flag.Parse()

	var (
		mp  = map[string]string{}
		err error
	)
	if *file != "" {
		fa, err := filepath.Glob(*file)
		if err != nil {
			log.Fatalln(err)
		}
		// 允许目录通配符,一个字符串匹配多个文件
		for _, v := range fa {
			mp[filepath.Base(v)] = v
		}
	} else {
		mp, err = getPathFromDB(*db, *title)
		if err != nil {
			log.Fatalln(err)
		}
	}

	sheet := make([]*xmind.Topic, 0, len(mp))
	for t, p := range mp {
		st, err := youDao(p, t)
		if err != nil {
			log.Fatalln(err)
		}
		sheet = append(sheet, st)
	}

	err = xmind.SaveSheets(*dst, sheet...)
	if err != nil {
		log.Fatalln(err)
	}
}

func youDao(src, title string) (*xmind.Topic, error) {
	fr, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer fr.Close()

	var node struct {
		Nodes json.RawMessage `json:"nodes"`
	}
	err = json.NewDecoder(fr).Decode(&node)
	if err != nil {
		return nil, err
	}

	// 有道云笔记思维导图,符合数组形式的结构,用自定义类型直接就可以转换
	st, err := xmind.LoadCustom([]byte(node.Nodes), map[string]string{
		xmind.CustomKeyId:       "id",
		xmind.CustomKeyTitle:    "topic",
		xmind.CustomKeyParentId: "parentid",
		xmind.CustomKeyIsRoot:   "isroot",
	})
	if err != nil {
		return nil, err
	}

	// 设置工作簿名称和中心主题名称,以及格式为逻辑图向右
	st.UpSheet(title, "", xmind.StructLogicRight)
	return st, nil
}

func getPathFromDB(p, title string) (map[string]string, error) {
	fileDir := filepath.Join(filepath.Dir(p), "file")
	fi, err := os.Stat(fileDir)
	if os.IsNotExist(err) {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not dir", fileDir)
	}

	fileMap := make(map[string]string) // 记录文件名和文件路径
	err = filepath.WalkDir(fileDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fileMap[d.Name()] = path
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", p)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer db.Close()

	row, err := db.Query("SELECT title,fileId FROM note WHERE title LIKE ?", title)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer row.Close()

	var (
		fileId string
		res    = make(map[string]string)
	)
	for row.Next() {
		err = row.Scan(&title, &fileId)
		if err != nil {
			return nil, err
		}

		// 根据数据库的id关联上磁盘的文件路径
		fp, ok := fileMap[fileId]
		if ok {
			res[title] = fp
		}
	}
	return res, row.Err()
}
