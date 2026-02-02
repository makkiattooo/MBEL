# クイックスタート: 最初のMBELアプリケーション

MBEL で本番対応の多言語アプリケーションを **15 分**で起動します。このガイドでは、プロジェクト設定、翻訳ファイル、コンパイル、デプロイメントについて説明します。

---

## 1. プロジェクトを初期化する

```bash
# プロジェクトディレクトリを作成
mkdir -p hello-mbel && cd hello-mbel

# Go モジュールを初期化
go mod init hello-mbel
go get github.com/makkiattooo/MBEL@latest

# プロジェクト構造を作成
mkdir -p locales cmd dist

# MBEL CLI をインストール
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

---

## 2. 翻訳ファイルを書く

### 日本語の翻訳を作成 (`locales/ja.mbel`)

```mbel
@namespace: hello
@lang: ja

app_name = "こんにちはMBEL"
app_version = "1.0.0"

greeting = "いらっしゃいませ、{name}さん！"
goodbye = "さようなら、{name}さん！ 時間：{time}"

# 日本語: [other] のみ
items_count(n) {
    [other] => "{n} 件のアイテムがあります"
}

profile_updated(gender) {
    [male]   => "彼がプロフィールを更新しました"
    [female] => "彼女がプロフィールを更新しました"
    [other]  => "彼らがプロフィールを更新しました"
}

ui.menu {
    home = "ホーム"
    about = "について"
    contact = "お問い合わせ"
    settings = "設定"
}

order_total = "合計: {price}（税金含む）"
```

---

## 3. 翻訳を検証してフォーマットする

```bash
# 構文エラーを確認
mbel lint locales/

# 自動フォーマット
mbel fmt locales/

# 統計を表示
mbel stats locales/
```

---

## 4. 翻訳をコンパイル

```bash
# JSON ファイルにコンパイル
mbel compile locales/ -o dist/translations.json

# ソースマップを含める
mbel compile locales/ -o dist/translations.json -sourcemap
```

---

## 5. Go アプリケーションを作成

### 基本的な例 (`cmd/main.go`)

```go
package main

import (
	"context"
	"fmt"
	"log"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

func main() {
	// 1. マネージャーを初期化
	m, err := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "ja",
		FallbackChain: []string{"ja", "en"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. 簡単な検索
	fmt.Println(m.Get("ja", "app_name", nil))

	// 3. 変数の補間
	vars := mbel.Vars{"name": "太郎"}
	greeting := m.Get("ja", "greeting", vars)
	fmt.Println(greeting)

	// 4. 複数形化
	for _, count := range []int{1, 5, 10} {
		vars := mbel.Vars{"n": count}
		msg := m.Get("ja", "items_count", vars)
		fmt.Printf("n=%d: %s\n", count, msg)
	}

	// 5. グローバル API
	mbel.Init(m)
	ctx := context.Background()
	fmt.Println(mbel.T(ctx, "greeting", mbel.Vars{"name": "花子"}))
}
```

---

## 6. テストを追加

```bash
go test -v ./...
```

---

## 7. ビルドして実行

```bash
go build -o hello-mbel ./cmd
./hello-mbel
```

---

## 8. 本番環境へのデプロイ

### オプション A: JSON を埋め込む

```go
//go:embed dist/translations.json
var translationsJSON []byte
```

### オプション B: Docker で配布

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

## 9. 次のステップ

1. **[マニュアル](Manual.md)** — 完全なドキュメント
2. **[ARCHITECTURE.md](ARCHITECTURE.md)** — 技術的な深い分析
3. **[DEVELOPMENT.md](DEVELOPMENT.md)** — MBEL を拡張
4. **[セキュリティベストプラクティス](SECURITY.md)** — XSS 対策

---

## トラブルシューティング

| 問題 | 解決策 |
|------|--------|
| `no such file or directory: locales/` | `mkdir -p locales` |
| `Undefined: mbel` | `go.mod` を確認 |
| `Syntax error at line 5` | `mbel lint locales/` |
| 翻訳が見つからない | キー、言語コード、フォールバック チェーンを確認 |
