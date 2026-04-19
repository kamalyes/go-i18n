# API 参考

## Manager

国际化管理器，负责消息加载、语言解析和翻译

### 构造函数

```go
NewManager(config *gci18n.I18N) (*Manager, error)
NewManagerWithLogger(config *gci18n.I18N, log logger.ILogger) (*Manager, error)
NewManagerWithFormatter(config *gci18n.I18N, formatter Formatter) (*Manager, error)
```

### 核心方法

```go
// 获取翻译消息（位置参数）
GetMessage(language, key string, args ...any) string

// 获取翻译消息（命名模板数据）
GetMessageWithMap(language, key string, templateData map[string]any) string

// 检查语言是否被支持
IsLanguageSupported(language string) bool

// 获取默认语言
GetDefaultLanguage() string

// 获取支持的语言列表
GetSupportedLanguages() []string

// 获取已加载的语言列表
GetLoadedLanguages() []string

// 重新加载指定语言的消息
ReloadLanguage(language string) error

// 解析语言代码
ResolveLanguage(lang string) string

// 检查指定语言的消息是否已加载
HasLanguage(language string) bool

// 获取指定语言的所有消息键
GetMessageKeys(language string) []string
```

## Context

国际化上下文，支持链式翻译操作

```go
ctx := &Context{Language: "en", Translator: manager}

// 翻译函数
ctx.T(key string, args ...any) string
ctx.TWithMap(key string, templateData map[string]any) string

// 获取/设置语言
ctx.GetLanguage() string
ctx.SetLanguage(language string)
```

## 全局翻译函数

```go
// 位置参数
T(ctx context.Context, key string, args ...any) string

// 命名模板数据
TWithMap(ctx context.Context, key string, templateData map[string]any) string

// 获取消息
GetMsgByKey(ctx context.Context, key string) string
GetMsgWithMap(ctx context.Context, key string, data map[string]any) string

// 获取/设置语言
GetLanguage(ctx context.Context) string
SetLanguage(ctx context.Context, language string) context.Context
```

## Loader

消息加载器接口

```go
type Loader interface {
    LoadMessages(language string) (map[string]string, error)
}
```

## Formatter

消息格式化器接口

```go
type Formatter interface {
    Format(message string, args []any, templateData map[string]any) string
}
```
