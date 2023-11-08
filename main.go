package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	name string
	path string
}

func NewFile(path string) *File {
	return &File{
		name: filepath.Base(path),
		path: path,
	}
}

var (
	root  string
	files []*File
)

func init() {
	flag.StringVar(&root, "path", "./AddonPackages", "AddonPackages path, eg: -path ./AddonPackages")
}

func main() {
	flag.Parse()
	os.Chdir(root)
	fmt.Println("开始遍历,当前路径为：", root)
	Walk("./")
	for _, file := range files {
		Work(file.path, file.name)
	}
	fmt.Println("程序运行结束")
}

func Walk(base string) {

	err := filepath.Walk(base, func(path string, info fs.FileInfo, err error) error {
		if path == "./" || path == "../" {
			return nil
		}
		if info.IsDir() {
			// Walk(strings.Join(stack, "/"))
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".var" {
			return nil
		}
		if info == nil {
			return nil
		}
		files = append(files, NewFile(path))
		// Work(path, info)
		return nil
	})
	if err != nil {
		fmt.Println("遍历错误：", err)
	}
}

func Work(path string, name string) {
	sep := strings.Split(name, ".")
	if len(sep) > 0 {
		author := sep[0]
		Check(author)
		fmt.Println("移动文件:", path, "->", filepath.Join(author, name))
		err := os.Rename(path, filepath.Join(author, name))
		if err != nil {
			fmt.Println("移动文件失败：", err)
		}
	}
}

func Check(path string) {
	_, err := os.Stat(path)
	if err != nil {
		// 文件夹不存在
		if os.IsNotExist(err) {
			fmt.Println("创建文件夹:", path)
			os.MkdirAll(path, os.ModePerm)
		}
	}
}
