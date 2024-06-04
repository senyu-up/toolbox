package file

import (
	"fmt"
	"testing"
)

func TestCreateFile(t *testing.T) {
	err := CreateFile("./test")
	fmt.Println(err)
}

func TestWebsite(t *testing.T) {

}

func TestDownloadFileToUrl(t *testing.T) {
	filePath, err := DownloadFileToUrl("https://s3.ap-east-1.amazonaws.com/test.download.file/20210525/20210525183823_%E4%BC%81%E4%B8%9A%E5%BE%AE%E4%BF%A120210525-183805%402x.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA2LZCRPKRKWIHSGFM%2F20210525%2Fap-east-1%2Fs3%2Faws4_request&X-Amz-Date=20210525T103826Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=2211a82814f7544faacadf2da306ea48ec9f04d4e27572a9bf37251e0d58cb10")
	fmt.Println(filePath, err)
}
