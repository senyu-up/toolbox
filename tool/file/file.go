package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

// 下载本地路径
const DownloadFilePath = "./running_time/upload/"

// 判断文件是否存在
func FileIsExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateFile(path string) error {
	if !FileIsExist(path) {
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
	}
	return nil
}

const LocalRoot = "./running_time/upload"

// 从网络上下载文件
// @Deprecated
func DownloadFileToUrl(fileUrl string) (filePath string, err error) {
	err = CreateFile(DownloadFilePath)
	if err != nil {
		return "", err
	}
	fileNameStr := path.Base(fileUrl)
	if fileNameStr == "" {
		return filePath, err
	}
	fileNameArr := strings.Split(fileNameStr, "?")
	fileName := fileNameArr[0]
	filePath = DownloadFilePath + fileName
	res, err := http.Get(fileUrl)
	if err != nil {
		fmt.Println("A error occurred!")
		return filePath, err
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	file, err := os.Create(filePath)
	if err != nil {
		return filePath, err
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d", written)
	return filePath, err
}

// ScanDirFile 递归获取指定目录下的所有文件
func ScanDirFile(dirPath string, includeChild bool, filter func(path, name string) bool) ([]string, error) {
	var result []string

	fis, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return result, err
	}

	// 所有文件/文件夹
	for _, fi := range fis {
		fullPathName := dirPath + "/" + fi.Name()
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		if filter != nil {
			if !filter(dirPath, fi.Name()) {
				continue
			}
		}
		// 是文件夹则递归进入获取;是文件，则压入数组
		if fi.IsDir() && includeChild {
			temp, err := ScanDirFile(fullPathName, includeChild, filter)
			if err != nil {
				return result, err
			}
			result = append(result, temp...)
			//	去掉隐藏文件
		} else {
			result = append(result, fullPathName)
		}
	}

	return result, nil
}

// IsExists 检查文件夹或文件是否存在
func IsExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// RemoveDir 移除传入的文件夹以及所有文件
func RemoveDir(dirPath string) error {
	return os.RemoveAll(dirPath)
}

func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}
