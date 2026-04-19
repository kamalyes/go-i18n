# go-i18n

功能强大的 Go 国际化（i18n）库，支持多种消息格式、语言回退和灵活的消息加载方式

## 特性

- **多种消息格式**：`fmt.Sprintf` 风格、`{key}`、`{{key}}`、`{{.key}}`
- **语言回退**：自动回退到默认语言
- **灵活加载**：Map、JSON、YAML、文件等多种来源
- **线程安全**：使用互斥锁保证并发安全
- **可扩展**：支持自定义格式化器和消息加载器

## 安装

```bash
go get github.com/kamalyes/go-i18n
```

## 基本用法

```go
loader := i18n.NewMapLoader(map[string]map[string]string{
    "en": {"hello": "Hello", "greeting": "Hello {name}"},
    "zh": {"hello": "你好", "greeting": "你好 {name}"},
})

config := &gci18n.I18N{
    Enable: true, DefaultLanguage: "en",
    SupportedLanguages: []string{"en", "zh"},
    MessageLoader: loader, EnableFallback: true,
}

manager, _ := i18n.NewManager(config)

manager.GetMessage("en", "hello")                                    // "Hello"
manager.GetMessage("zh", "hello")                                    // "你好"
manager.GetMessage("en", "greeting", "John")                        // "Hello John"
manager.GetMessageWithMap("en", "greeting", map[string]any{"name": "John"}) // "Hello John"
```

## 框架图

```mermaid
graph TB
    subgraph Core
        M[Manager]
        F[Formatter]
        T[Translator]
        C[Context]
    end
    subgraph Loaders
        L[Loader]
        ML[MapLoader]
        FL[FileLoader]
        JL[JSONLoader]
        YL[YAMLLoader]
    end
    subgraph Utils
        TF[T/TWithMap]
    end
    M --> T
    M --> F
    C --> T
    TF --> C
    L --> M
    ML --> L
    FL --> L
    JL --> L
    YL --> L
```

## 核心组件

| 组件 | 文件 | 说明 |
|------|------|------|
| Manager | manager.go | 国际化管理器，负责消息加载、语言解析和翻译 |
| Formatter | format.go | 消息格式化器，支持多种占位符格式 |
| Translator | context.go | 翻译器接口，解耦 Context 与 Manager |
| Context | context.go | 国际化上下文，支持链式翻译操作 |
| Loader | loader.go | 消息加载器接口 |
| MapLoader | loader_map.go | 内存 Map 加载器 |
| FileLoader | loader_file.go | JSON 文件加载器 |
| JSONLoader | loader_json.go | JSON 字符串加载器 |
| YAMLLoader | loader_yaml.go | YAML 文件/字符串加载器 |

## 测试

```bash
# 测试所有组件并生成覆盖率报告
go test -cover ./... -v
```

## 文档

| 文档 | 说明 |
|------|------|
| [docs/API.md](docs/API.md) | API 参考 |
| [docs/Loaders.md](docs/Loaders.md) | 消息加载器 |
