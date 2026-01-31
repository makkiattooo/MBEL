# MBEL: 公式完全リファレンスマニュアル

**バージョン:** 1.2.0
**日付:** 2026年1月

---

## 📖 目次

1.  [はじめに](#1-はじめに)
2.  [MBEL 言語](#2-mbel-言語)
    *   [ファイル構造](#21-ファイル構造)
    *   [データ型](#22-データ型)
    *   [補間と変数](#23-補間と変数)
    *   [ロジックと制御フロー](#24-ロジックと制御フロー)
3.  [CLI ツールチェーン](#3-cli-ツールチェーン)
4.  [Go SDK の統合](#4-go-sdk-の統合)

---

## 1. はじめに

Modern Billed-English Language (MBEL) は、AI時代のために設計された次世代のローカリゼーションシステムです。

---

## 2. MBEL 言語

### 2.1 ファイル構造

```mbel
@namespace: "features.auth"  # メタデータ

[ログイン画面]               # セクション

title = "ログイン"           # 代入
```

### 2.2 データ型

*   **文字列リテラル**: 二重引用符で囲みます。
*   **複数行文字列**: 三重引用符 `"""` で囲みます。

### 2.3 補間と変数

変数は波括弧 `{}` で囲みます。

*   **構文**: `こんにちは、{user_name}さん！`

### 2.4 ロジックと制御フロー

**構文:** `key(variable) { cases }`

#### 完全一致
```mbel
theme(mode) {
    [dark]  => "ダークモード"
    [light] => "ライトモード"
}
```

#### 範囲一致
```mbel
battery(percent) {
    [0]       => "空"
    [1-19]    => "バッテリー残量低下"
    [100]     => "満充電"
}
```

#### 複数形
日本語には文法的な複数形変化はありませんが、ロジックブロックを使用して数値による条件分岐を行うことができます。通常、`[other]` ルールのみを使用します。

```mbel
items_ja(n) {
    [other] => "{n} 個のアイテム"
}
```

---

## 3. なぜ MBEL なのか？（実際の比較）

まだ JSON で十分だと思っていませんか？ 単純なロジック + 補間を比較してみましょう。

#### ❌ JSON 方式（煩雑）
```json
{
  "cart_items_other": "カートに {{count}} 個の商品があります。",
  "greeting_male": "お帰りなさい、{{name}} 様",
  "greeting_female": "お帰りなさい、{{name}} 様"
}
```
*ロジックが複数のキーに分散しています。Go/JS コードはどのキーを取得するかを判断する必要があります。*

#### ✅ MBEL 方式（クリーン）
```mbel
cart_items(n) {
    [other] => "カートに {n} 個の商品があります。"
}

greeting(gender) {
    [male]   => "お帰りなさい、{name} 様"
    [female] => "お帰りなさい、{name} 様"
}
```
*1つのキー、クリーンなロジック。ランタイムが複雑な処理をハンドルします。*

---

## 4. 構文ガイド

### 4.1 基本キー
単純なキーと値のペアです。

```mbel
title = "私のアプリケーション"
```

### 4.2 補間 vs ロジック変数
重要な違い：
1.  **制御変数**: `key(var)` 内の変数。どのケースを選択するかを決定します。
2.  **補間変数**: `{var}` 内の変数。単にテキストとして置き換えられます。

```mbel
# 'gender' は制御変数
# '{name}' は補間変数
greeting(gender) {
    [male]   => "こんにちは、{name} 様"
    [female] => "こんにちは、{name} 様"
}
```
*実行時の使用方法:* `mbel.T(ctx, "greeting", mbel.Vars{"gender": "male", "name": "田中"})`

### 4.3 AI メタデータ
メタデータはコンパイルされたオブジェクトの `__ai` フィールドに格納されます。実行時のテキストには影響しませんが、翻訳エージェントの能力を大幅に強化します。

---

## 5. CLI ツールチェーン

*   `mbel init`: セットアップウィザード。
*   `mbel lint`: 構文チェック。
*   `mbel compile`: JSONへのコンパイル。
*   `mbel watch`: ホットリロード（開発用）。
*   **`stats`**: 統計情報。
*   **`fmt`**: コードフォーマット。

---

## 4. Go SDK の統合

### 4.1 アーキテクチャ

*   **Manager**: 中央エントリーポイント。
*   **Runtime**: 実行環境。
*   **Repository**: データソースインターフェース。

### 4.2 初期化

```go
import "github.com/makkiattooo/MBEL"

func init() {
    mbel.Init("./locales", mbel.Config{
        DefaultLocale: "ja",
        Watch:         true,
    })
}
```

### 4.3 使用法 (T 関数)

`T` (Translate) 関数はコンテキストに基づいて文字列を解決します。

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // 1. シンプルなキー
    title := mbel.T(ctx, "page_title")

    // 2. 変数付き
    msg := mbel.T(ctx, "welcome", mbel.Vars{"name": "佐藤"})

    // 3. ロジック/複数形
    items := mbel.T(ctx, "cart_items", 5) 
}
```

### 4.4 HTTP ミドルウェア

MBELは `Accept-Language` ヘッダーを自動的に解析します。

```go
router.Use(mbel.Middleware)
```
