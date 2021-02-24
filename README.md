### Asura

> Blazing fast http pluginHook for humans

##### 项目简介

- Restful Router及单元测试的良好支持
- 提供类似于 Laravel 的 middleware(Filters & Terminators) 机制
- 提供了统一的Exception处理层
- 设计了Validator/Binding帮助我们更好构建Request和转换数据类型
- 灵活的适配器设计模式，统一了Response数据格式

##### 项目特点

- 模块化设计，核心足够轻量

##### 编译环境

- **请只用 Golang v1.11.x 以上版本编译执行**

##### 如何使用

```go
package main

import (
	Asura "github.com/Shies/Asura"
)

func main() {
    engine := Asura.Default()
    engine.GET("/ping", func(c *Asura.Context) {
    	c.String(200, "%s", "pong")
    })
    engine.Run(":8080")
}
```

##### 依赖包

- Go Modules dependency
