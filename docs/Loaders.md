# 消息加载器

## Loader 接口

```go
type Loader interface {
    LoadMessages(language string) (map[string]string, error)
}
```

## MapLoader

从内存 Map 加载消息，适用于程序启动时硬编码或动态构建的翻译数据

```go
loader := i18n.NewMapLoader(map[string]map[string]string{
    "en": {"hello": "Hello", "greeting": "Hello {name}"},
    "zh": {"hello": "你好", "greeting": "你好 {name}"},
})
```

## FileLoader

从 JSON 文件加载消息

```go
loader := i18n.NewFileLoader("path/to/locales")
// 文件: path/to/locales/en.json
```

JSON 文件格式：

```json
{
    "hello": "Hello",
    "greeting": "Hello {name}"
}
```

## JSONLoader

从 JSON 字符串加载消息

```go
loader, _ := i18n.NewJSONLoader(`{"en": {"hello": "Hello"}}`)
```

## YAMLLoader

从 YAML 文件加载消息

```go
loader := i18n.NewYAMLLoader("path/to/locales")
// 文件: path/to/locales/en.yaml
```

## YAMLStringLoader

从 YAML 字符串加载消息

```go
loader, _ := i18n.NewYAMLStringLoader(`
en:
  hello: Hello
zh:
  hello: 你好
`)
```

## 嵌套结构扁平化

所有加载器支持嵌套 JSON/YAML 自动扁平化为点号格式：

```json
{
    "error": {
        "internal": "Internal Error",
        "external": "External Error"
    }
}
```

自动转换为：

```
error.internal: Internal Error
error.external: External Error
```
