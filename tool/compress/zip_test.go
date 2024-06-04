package compress

import (
	"os"
	"path"
	"runtime"
	"testing"
)

func chToRunCwd() {
	// 获取当前文件运行的目录
	_, filename, _, _ := runtime.Caller(0)
	os.Chdir(path.Dir(filename))
}

func SetUp() {
	chToRunCwd()
	// 获取当前文件运行的目录
	os.Mkdir("./empty_dir", os.ModePerm)
	// go touch file
	os.WriteFile("./no_data.zip", []byte("0"), os.ModePerm)
}

func TearDown() {
	chToRunCwd()
	os.RemoveAll("./empty_dir")
	os.Remove("./no_data.zip")

	os.Remove("./example.zip")
	os.Remove("./empty_dir.zip")
	os.Remove("./test_file.zip")
	os.Remove("./test_file2.zip")

	os.Remove("./go.mod")
}

func TestZip(t *testing.T) {
	SetUp()
	type args struct {
		srcPath string
		destZip string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1", args: args{srcPath: "./main/", destZip: "./main.zip"}, wantErr: true, // 文件夹不存在
		},
		{
			name: "2", args: args{srcPath: "./example/", destZip: "./example.zip"}, wantErr: false, // 文件夹存在
		},
		{
			name: "3", args: args{srcPath: "./empty_dir/", destZip: "./empty_dir.zip"}, wantErr: false, // 空文件夹
		},
		{
			name: "4", args: args{srcPath: "./zip_test.go", destZip: "./test_file.zip"}, wantErr: false, // 文件
		},
		{
			name: "5", args: args{srcPath: "./zip_test2.go", destZip: "./test_file2.zip"}, wantErr: true, // 文件
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Zip(tt.args.srcPath, tt.args.destZip); (err != nil) != tt.wantErr {
				t.Errorf("Zip() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				// 不报错就 检查文件是否存在
				if _, err := os.Stat(tt.args.destZip); os.IsNotExist(err) {
					t.Errorf("Zip() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}

	TearDown()
}

func TestUnZip(t *testing.T) {
	SetUp()
	type args struct {
		srcFile string
		destDir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1", args: args{srcFile: "./go_mod.zip", destDir: "./"}, wantErr: false, // 文件存在
		},
		{
			name: "2", args: args{srcFile: "./aaa_bbb.zip", destDir: "./"}, wantErr: true, // 文件不存在
		},
		{
			name: "3", args: args{srcFile: "./main.zip", destDir: "./"}, wantErr: false, // 空的 zip
		},
		{
			name: "4", args: args{srcFile: "./no_data.zip", destDir: "./"}, wantErr: true, // 捏造的 zip
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnZip(tt.args.srcFile, tt.args.destDir); (err != nil) != tt.wantErr {
				t.Errorf("UnZip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	TearDown()
}
