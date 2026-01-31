# MBEL：技巧与最佳实践

遵循这些模式以充分发挥 MBEL 的作用。

## 1. 命名约定
使用分层键，而不是像 `label1` 这样的扁平名称。
*   **模式**：`[feature].[screen].[component].[element]`

## 2. 文件组织
将您的本地化划分为逻辑单元。文件夹结构通过 `--ns` 选项自动定义命名空间。

## 3. 使用 AI 上下文
不要吝啬使用 `@AI_Context`。它是自动化翻译的质量保证。

## 4. 逻辑块 vs 代码中的插值
避免在 Go 代码中构建句子。
*   **错误 (Go)**：`fmt.Sprintf(mbel.T(ctx, "hello") + " " + name)`
*   **正确 (MBEL)**：
    ```mbel
    welcome(name) {
        [other] => "欢迎，{name}！"
    }
    ```
