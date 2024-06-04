package file

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func ScanDirByName(basePath string, maxDepth int, name string) (list []string, err error) {
	return ScanDir(basePath, maxDepth, &Pattern{Name: name})
}

func ScanFileByName(basePath string, maxDepth int, name string) (list []string, err error) {
	return ScanFile(basePath, maxDepth, &Pattern{Name: name})
}

func ScanFile(filePath string, maxDepth int, pattern *Pattern) (list []string, err error) {
	defer func() {
		if list == nil || len(list) == 0 {
			err = errors.New("file not find by pattern")
		}
	}()
	list = make([]string, 0, 1)
	err = doScan(filePath, 0, maxDepth, false, pattern, &list)

	return
}

func ScanDir(filePath string, maxDepth int, pattern *Pattern) (list []string, err error) {
	defer func() {
		if list == nil || len(list) == 0 {
			err = errors.New("file not find by pattern")
		}
	}()
	list = make([]string, 0, 1)
	err = doScan(filePath, 0, maxDepth, true, pattern, &list)

	return
}

func doScan(dir string, curDepth int, maxDepth int, scanDir bool, pattern *Pattern, list *[]string) (err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil {
			return err
		}

		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		if (info.IsDir() && scanDir) || (!info.IsDir() && !scanDir) {
			if pattern.Name != "" {
				if pattern.Name == info.Name() {
					*list = append(*list, path)
				}
			} else if pattern.Pattern != nil {
				if pattern.Pattern.MatchString(info.Name()) {
					*list = append(*list, path)
				}
			}
		}

		return nil
	})

	return err
}

type Pattern struct {
	// 匹配特定的文件名, 与pattern二选一, 如果两个参数都赋值, 则优先使用name
	Name string
	// 基于正则查询符合的文件, 与 pattern 二选一
	Pattern *regexp.Regexp
}

// ScanConfigPath
//
//	@Description: 在这个目录下搜寻config文件
//	@param basePath  body any true "-"
//	@return string
func ScanConfigPath(basePath string) string {
	confPath, err := ScanDirByName(basePath, 10, "config")
	if err != nil {
		return ""
	} else if len(confPath) < 1 {
		return ""
	}
	return path.Join(confPath[0], "config.yaml")
}
