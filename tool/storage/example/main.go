package main

import (
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/storage"
	"os"
	"time"
)

func main() {
	var conf = &config.Aws{
		AwsAccessId:  "123",
		AwsAccessKey: "xxx",

		S3: []config.S3Storage{
			{
				Region: storage.AWSRegion,
				Bucket: "test",
				Path:   "xh_test/gm",
				Expire: 24,
				Host:   "https://cdn.sug.com/images/cd1",
			},
			{
				Region: storage.HKAWSRegion,
				Bucket: "test2",
				Path:   "xh_test2/gm",
				Expire: 24,
				Host:   "https://cdn.sug.com/images/hk1",
			},
			{
				Region: storage.AWSUwRegion,
				Bucket: "risk.resource.test",
				Path:   "xh_test2/gm",
				Expire: 24,
				Host:   "https://cdn.sug.com/images/hk1",
			},
		},
	}
	s3Client := storage.InitByConf(conf)

	// 上传文件
	f, _ := os.Open("test.jpg")
	if _, err := s3Client.UploadFileCDN(f, "test.jpg", storage.AWSUwRegion); err != nil {
		fmt.Printf("upload err %v \n", err)
	} // 这个地区没初始化，所以会失败

	// 上传文件
	if _, err := s3Client.UploadFileCDN(f, "test.jpg", storage.HKAWSRegion); err != nil {
		fmt.Printf("upload 2 err %v \n", err)
	}

	// 获取加密文件地址
	if url, cdn, err := s3Client.GetPreSignUrlCdn(storage.AWSUwRegion, "test.jpg", "", ""); err != nil {
		fmt.Printf("get url err %v \n", err)
	} else {
		fmt.Printf("cdn %s url %s \n", cdn, url)
	}

	// 获取加密文件地址
	var imgPath = "/develop/NTRYbVhiallhYToxNjc2MDE4MTMyOmRldmVsb3A=/v1/1b791790-fd70-4dc3-baa9-ded08dc143a4_202305151626887.png"
	if url, err := s3Client.GetPreSignUrl(time.Hour, imgPath, storage.AWSUwRegion); err != nil {
		fmt.Printf("get url err %v \n", err)
	} else {
		fmt.Printf("cdn %s \n", url)
	}
}
