package main

import (
	"fmt"
	"github.com/senyu-up/toolbox/tool/file"
	"regexp"
)

func Example1() {
	// 通过 基础路径，和名字 查找文件夹
	dirs, err := file.ScanDirByName("/Users/th/xh/github.com/senyu-up/toolbox", 5, "example")
	if err != nil {
		fmt.Printf("scan dir err %v \n", err)
		return
	} else {
		fmt.Printf("get scan dir re %v \n", dirs)
	}

	// 通过 基础路径，和名字 查找文件夹
	reg, err := regexp.Compile("\\.md")
	if err != nil {
		fmt.Printf("new regex err %v\n", err)
		return
	}
	files, err := file.ScanFile("/Users/th/xh/toolbox", 5, &file.Pattern{Pattern: reg})
	if err != nil {
		fmt.Printf("scan  files err %v \n", err)
		return
	} else {
		fmt.Printf("get scan file re %v \n", files)
	}
}

func Example2() {
	if p, err := file.LookupFile("/Users/th/xh/toolbox/tool/file/example", "go.mod", 5); err != nil {
		fmt.Printf("look up file err %v \n", err)
	} else {
		fmt.Printf("look up file: %s \n", p)
	}

	if d, err := file.LookupDir("/Users/th/xh/toolbox/tool/file/example", "config", 5); err != nil {
		fmt.Printf("look up file err %v \n", err)
	} else {
		fmt.Printf("look up file: %s \n", d)
	}
}

func Example4() {
	file.IsDir("./toolbox/tool/file/example") // 检查 文件夹 是否存在

	file.IsExists("./toolbox/tool/file/example") // 检查 文件、文件夹 是否存在

	file.FileIsExist("./toolbox/tool/file/example/main.go") // 检查 文件 是否存在

}

func main() {
	Example1()
	//Example2()

	Example4()
}
