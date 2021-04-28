### MicroServe

> Blazing fast UDP & Gateway & DTS server for humans

##### 核心简介

- Restful Router及单元测试的良好支持
- 提供类似于 Laravel 的 middleware(Filters & Terminators) 机制
- 提供了统一的Exception处理层
- 设计了Validator/Binding帮助我们更好构建Request和转换数据类型
- 灵活的适配器设计模式，统一了Response数据格式

##### 项目特点

- 模块化设计，核心足够轻量

##### 项目介绍

- UDP传输服务器
- API网关服务器
- ALIYUN-DTS-缓存订阅中心
- ...

##### 编译环境

- **请只用 Golang v1.13.x 以上版本编译执行**
- **linux supervisor监控进程**
- **cmd/bin启动服务**

##### 如何使用

```shell script
export GO111MODULE="on"
cd app/cmd
GOOS="linux" go build -o ./udpserver
./udpserver -conf config.example.toml
```

```shell script
export GO111MODULE="on"
cd app/cmd
GOOS="linux" go build -o ./gwserver
./gwserver -conf config.example.toml
```

```shell script
export GO111MODULE="on"
cd app/cmd
go run main.go -conf config.example.toml
# go run main.go &
```

```shell script
export GO111MODULE="on"
cd app/cmd
GOOS="linux" go install
$GOPATH/cmd/cmd -conf config.example.toml
```

##### 依赖包

- Go Modules dependency

##### 特技

- 使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
- Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
- 你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
- [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
- Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
- Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
