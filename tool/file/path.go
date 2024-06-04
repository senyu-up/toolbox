package file

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

func GetAbsPath(path string) string {
	if !filepath.IsAbs(path) {
		abs, _ := filepath.Abs(path)
		return abs
	}

	return path
}

// LookUpwardFile
// @description 向上递归查找指定文件
func LookUpwardFile(dir string, filename string, depth int, excludes ...string) (filepath string, err error) {
	dir = GetAbsPath(dir)
	var excludeMap map[string]struct{}
	if excludes != nil {
		excludeMap = make(map[string]struct{}, len(excludeMap))
		for _, tmpDir := range excludes {
			excludeMap[GetAbsPath(tmpDir)] = struct{}{}
		}
	}

	return
}

type pathType int

const (
	pathTypeFile pathType = 1
	pathTypeDir  pathType = 2
)

var NotFoundError = errors.New("file not found")

// LookupFile
// @description 向上寻找目标文件
// depth 为0时 无限向上查找, 直到查询到"/"(不包含)为止
// depth = 1 时, 只查找当前目录
func LookupFile(dir string, name string, depth int) (filepath string, err error) {
	dir = GetAbsPath(dir)

	return doLookUp(dir, name, pathTypeFile, 1, depth)
}

// LookupDir
// @description 从指定的dir目录, 向上寻找目标目录
// depth 为0时, 无限向上查找, 直到查询到"/"(不包含)为止
// depth = 1 时, 只查找当前目录
func LookupDir(dir string, name string, depth int) (dirPath string, err error) {
	dir = GetAbsPath(dir)
	return doLookUp(dir, name, pathTypeDir, 1, depth)
}

func doLookUp(dir string, name string, pt pathType, curDepth int, depth int) (filepath string, err error) {
	// 不允许扫描根目录
	if dir == "/" || (depth > 0 && curDepth > depth) {
		return "", NotFoundError
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	var isDir bool
	for i, _ := range files {
		if strings.HasPrefix(files[i].Name(), ".") {
			continue
		}
		if isDir = files[i].IsDir(); isDir {
			if pt == pathTypeDir && files[i].Name() == name {
				return dir + "/" + files[i].Name(), nil
			}
		} else {
			if pt == pathTypeFile && files[i].Name() == name {
				return dir + "/" + files[i].Name(), nil
			}
		}
	}

	return doLookUp(path.Dir(dir), name, pt, curDepth+1, depth)
}
