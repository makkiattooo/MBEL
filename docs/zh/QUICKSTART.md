# 快速入门：你的第一个 MBEL 应用程序

使用 MBEL 在 **15 分钟**内启动生产就绪的多语言应用程序。本指南涵盖项目设置、翻译文件、编译和部署。

---

## 1. 初始化项目

```bash
# 创建项目目录
mkdir -p hello-mbel && cd hello-mbel

# 初始化 Go 模块
go mod init hello-mbel
go get github.com/makkiattooo/MBEL@latest

# 创建项目结构
mkdir -p locales cmd dist

# 安装 MBEL CLI
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

---

## 2. 编写翻译文件

### 创建中文翻译 (`locales/zh.mbel`)

```mbel
@namespace: hello
@lang: zh

app_name = "你好 MBEL"
app_version = "1.0.0"

greeting = "欢迎，{name}！"
goodbye = "再见，{name}！时间：{time}"

# 中文: [other] 形式
items_count(n) {
    [other] => "你有 {n} 个项目"
}

profile_updated(gender) {
    [male]   => "他更新了他的个人资料"
    [female] => "她更新了她的个人资料"
    [other]  => "他们更新了他们的个人资料"
}

ui.menu {
    home = "主页"
    about = "关于"
    contact = "联系"
    settings = "设置"
}

order_total = "总计：{price}（含税）"
```

---

## 3. 验证和格式化翻译

```bash
# 检查语法错误
mbel lint locales/

# 自动格式化
mbel fmt locales/

# 显示统计信息
mbel stats locales/
```

---

## 4. 编译翻译

```bash
# 编译为单个 JSON 文件
mbel compile locales/ -o dist/translations.json

# 包含源映射以供调试
mbel compile locales/ -o dist/translations.json -sourcemap
```

---

## 5. 创建 Go 应用程序

### 基本示例 (`cmd/main.go`)

```go
package main

import (
	"context"
	"fmt"
	"log"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

func main() {
	// 1. 初始化 Manager
	m, err := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "zh",
		FallbackChain: []string{"zh", "en"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. 简单查询
	fmt.Println(m.Get("zh", "app_name", nil))

	// 3. 变量插值
	vars := mbel.Vars{"name": "李明"}
	greeting := m.Get("zh", "greeting", vars)
	fmt.Println(greeting)

	// 4. 复数形式
	for _, count := range []int{1, 5, 10} {
		vars := mbel.Vars{"n": count}
		msg := m.Get("zh", "items_count", vars)
		fmt.Printf("n=%d: %s\n", count, msg)
	}

	// 5. 全局 API
	mbel.Init(m)
	ctx := context.Background()
	fmt.Println(mbel.T(ctx, "greeting", mbel.Vars{"name": "王芳"}))
}
```

---

## 6. 添加测试

```bash
go test -v ./...
```

---

## 7. 构建并运行

```bash
go build -o hello-mbel ./cmd
./hello-mbel
```

---

## 8. 生产部署

### 选项 A：嵌入 JSON

```go
//go:embed dist/translations.json
var translationsJSON []byte
```

### 选项 B：使用 Docker 分发

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /src
COPY . .
RUN go build -o /tmp/hello-mbel ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /tmp/hello-mbel .
COPY dist/translations.json ./dist/
CMD ["./hello-mbel"]
```

---

## 9. 后续步骤

1. **[手册](Manual.md)** — 完整文档
2. **[ARCHITECTURE.md](ARCHITECTURE.md)** — 深度技术分析
3. **[DEVELOPMENT.md](DEVELOPMENT.md)** — 扩展 MBEL
4. **[安全最佳实践](SECURITY.md)** — XSS 防护

---

## 故障排除

| 问题 | 解决方案 |
|------|--------|
| `no such file or directory: locales/` | `mkdir -p locales` |
| `Undefined: mbel` | 检查 `go.mod` |
| `Syntax error at line 5` | `mbel lint locales/` |
| 未找到翻译 | 检查密钥、语言代码和备用链 |
