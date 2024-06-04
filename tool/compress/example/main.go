package main

import (
	"encoding/base64"
	"fmt"
	"github.com/senyu-up/toolbox/tool/compress"
)

func Example3() {
	// 压缩一个文件
	if err := compress.Zip("./go.mod", "./go_mod.zip"); err != nil {
		fmt.Printf("zip file err %v \n", err)
	} else {
		fmt.Printf("zip file success \n")
	}
	// 压缩一个文件夹
	if err := compress.Zip("./example", "./example.zip"); err != nil {
		fmt.Printf("zip dir err %v \n", err)
	} else {
		fmt.Printf("zip dir success \n")
	}
	// 解压文件
	if err := compress.UnZip("./main.zip", "./aaa/"); err != nil {
		fmt.Printf("unzip file err %v \n", err)
	} else {
		fmt.Printf("unzip file success \n")
	}
}

func Example4() {
	// gzip 压缩
	if data, err := compress.Gzip([]byte("hello world")); err != nil {
		fmt.Printf("gzip err %v \n", err)
	} else {
		fmt.Printf("gzip success %v \n", data)
	}
	data, _ := base64.StdEncoding.DecodeString("H4sIAAAAAAAA/8pIzcnJVyjPL8pJUQQAlQYlBQAAAA==")
	// gzip 解压缩
	if data, err := compress.Gunzip(data); err != nil {
		fmt.Printf("gunzip err %v \n", err)
	} else {
		fmt.Printf("gunzip success %v \n", string(data))
	}
}

func main() {
	//Example3()

	Example4()
}
