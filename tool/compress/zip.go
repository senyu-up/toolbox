package compress

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Zip
//
//	@Description: 把指定文件夹/文件 压缩到 指定文件。 地址不合法，文件不存在 会导致报错
//	@param srcFile  body any true "-"
//	@param destZip  body any true "-"
//	@return error
func Zip(srcPath string, destZip string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	err = filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(srcPath)+"/")
		// header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})
	return err
}

// UnZip 解压zip文件到指定目录
func UnZip(srcFile string, destDir string) error {
	zipReader, err := zip.OpenReader(srcFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	//遍历reader 获取所有文件和目录
	for _, f := range zipReader.File {
		path := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return err
			}
			//获取文件reader
			fr, err := f.Open()
			if err != nil {
				return err
			}
			//创建写文件的writer
			fw, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(fw, fr)
			if err != nil {
				return err
			}
			fr.Close()
			fw.Close()
		}
	}
	return nil
}
