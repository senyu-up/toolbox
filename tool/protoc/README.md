# protoc工具安装说明

1. 进入本项目的tool/protoc文件

2. sudo sh install.sh (必须使用zsh)

3. 成功后重启当前控制台，输入

4. protoc --version

5. 若打印出

6. libprotoc 3.14.0 则表示成功安装protoc

7. 安装插件

8. go get -u github.com/golang/protobuf/protoc-gen-go

9. go get -u github.com/gogo/protobuf/protoc-gen-gogofast

10. go get github.com/gogo/protobuf/protoc-gen-gogofaster

11. go get github.com/gogo/protobuf/protoc-gen-gogoslick

    > 若安装完gogo proto插件后，使用shell脚本无法生成pb文件时且报如下错误：
    >
    > ```
    > github.com/gogo/protobuf/gogoproto/gogo.proto: File not found.
    > common.proto:6:1: Import "github.com/gogo/protobuf/gogoproto/gogo.proto" was not found or had errors.
    > ```
    >
    > a.可能是你在非toolbox工程里生成pb文件，需要把toolbox下的**<u>lib</u>**文件夹拷贝到你的工程中。

# gogoproto如何生成pb?

- 进入proto文件夹
- sh gogo.sh
- 输入要生成pb的文件夹名称
- 回车

# 如何让goland识别gogoproto import

1. Preferences->
2. Language&Frameworks->
3. Protocol Buffers
4. import paths 里点+，添加本项目的proto文件夹，即{yourpath}/toolbox/lib


# gogoprotobuf示例
```
syntax = "proto3";

import "github.com/gogo/protobuf/gogoproto/gogo.proto";

//package名称和 文件名 一致
package order;

//go_package名称和 文件夹 名称一致
option go_package = "gogo_demo";

service Order{
// 获取订单详情
rpc GetOrderDetail(GetOrderDetailReq)returns(GetOrderDetailResp);
}

message GetOrderDetailReq {
// 订单号
string trans_no = 1 [(gogoproto.moretags) = 'form:"trans_no" validate:"required"'];
}

message GetOrderDetailResp {
// 订单ID
uint64 id = 1 [(gogoproto.jsontag) = 'id'];
// 订单编号
string trans_no = 2 [(gogoproto.jsontag) = 'trans_no'];
// 下单时间
string created_at = 3 [(gogoproto.jsontag) = 'created_at'];
// 订单状态
int32 search_status = 4 [(gogoproto.jsontag) = 'search_status'];
}
```

# proto文件引暂不可跨文件夹相互import

# 安装oh my zsh
- 首先确定你当前的shell是 zsh
- `echo $SHELL` 输出应为 `/bin/zsh`
- 若不是zsh,使用命令`chsh -s /bin/zsh`切换到zsh
- 使用如下命令全自动安装oh my zsh
  `sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"`
- 若失败（资源被墙等问题）
- 源码安装方案：
```
  git clone https://github.com/ohmyzsh/ohmyzsh.git ~/.oh-my-zsh
  echo "source ~/.oh-my-zsh/oh-my-zsh.sh" >> ~/.zshrc
  git clone --depth=1 https://gitee.com/romkatv/powerlevel10k.git ~/powerlevel10k
  echo 'source ~/powerlevel10k/powerlevel10k.zsh-theme' >>~/.zshrc
  source ~/.zshrc
  此时会有一些config要做,按照引导输入就行
```
# 安装oh my zsh 插件
```
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting
git clone https://github.com/zsh-users/zsh-autosuggestions ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions

sed -i '' '1i\
plugins=(git zsh-syntax-highlighting zsh-autosuggestions)
' ~/.zshrc

source ~/.zshrc
```