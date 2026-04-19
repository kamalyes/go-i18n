/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-29 00:00:00
 * @FilePath: \go-i18n\loader_yaml.go
 * @Description: YAML 消息加载器，从 YAML 文件或字符串加载翻译消息
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"fmt"
	"os"
	"path/filepath"

	gci18n "github.com/kamalyes/go-config/pkg/i18n"
	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/errorx"
	"github.com/kamalyes/go-toolbox/pkg/json"
)

// 确保 YAMLLoader 同时实现 Loader 和 MessageLoader 接口
var (
	_ Loader               = (*YAMLLoader)(nil)
	_ gci18n.MessageLoader = (*YAMLLoader)(nil)
)

// YAMLLoader YAML 文件消息加载器
// 从指定目录下的 YAML 文件加载翻译消息，支持嵌套 YAML 自动扁平化
// 文件命名规则: {localesPath}/{language}.yaml 或 {localesPath}/{language}.yml
type YAMLLoader struct {
	// localesPath 翻译文件所在目录路径
	localesPath string
}

// NewYAMLLoader 创建 YAML 文件消息加载器
// localesPath: 翻译文件所在目录路径，如 "./locales" 或 "resources/locales"
func NewYAMLLoader(localesPath string) *YAMLLoader {
	return &YAMLLoader{
		localesPath: localesPath,
	}
}

// LoadMessages 加载指定语言的消息
// 依次尝试 {language}.yaml 和 {language}.yml 文件
// 支持嵌套 YAML 自动扁平化为点号格式
// language: 语言代码，如 "zh", "en", "ja" 等
func (y *YAMLLoader) LoadMessages(language string) (map[string]string, error) {
	yamlPath := filepath.Join(y.localesPath, language+".yaml")
	ymlPath := filepath.Join(y.localesPath, language+".yml")

	filePath := ""
	if _, err := os.Stat(yamlPath); err == nil {
		filePath = yamlPath
	} else if _, err := os.Stat(ymlPath); err == nil {
		filePath = ymlPath
	}

	if filePath == "" {
		return nil, errorx.NewError(ErrTypeLoaderFileNotFound, fmt.Sprintf("路径: %s 或 %s", yamlPath, ymlPath))
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errorx.NewError(ErrTypeLoaderReadFailed, fmt.Sprintf("文件: %s, 错误: %v", filePath, err))
	}

	jsonData, err := convert.YAMLToJSON(data)
	if err != nil {
		return nil, errorx.NewError(ErrTypeLoaderParseFailed, fmt.Sprintf("文件: %s, YAML 转 JSON 失败: %v", filePath, err))
	}

	var nested map[string]any
	if err := json.Unmarshal(jsonData, &nested); err != nil {
		return nil, errorx.NewError(ErrTypeLoaderParseFailed, fmt.Sprintf("文件: %s, JSON 解析失败: %v", filePath, err))
	}

	return FlattenToMessages(nested), nil
}

// 确保 YAMLStringLoader 同时实现 Loader 和 MessageLoader 接口
var (
	_ Loader               = (*YAMLStringLoader)(nil)
	_ gci18n.MessageLoader = (*YAMLStringLoader)(nil)
)

// YAMLStringLoader YAML 字符串消息加载器
// 从 YAML 格式字符串直接解析翻译消息，适用于配置文件内嵌或动态生成的场景
// YAML 格式为: zh: key: 值  en: key: value
type YAMLStringLoader struct {
	// messages 已解析的语言消息映射
	messages map[string]map[string]string
}

// NewYAMLStringLoader 创建 YAML 字符串消息加载器
// yamlStr: YAML 格式的翻译消息字符串
// 格式示例:
//
//	zh:
//	  hello: 你好
//	en:
//	  hello: Hello
func NewYAMLStringLoader(yamlStr string) (*YAMLStringLoader, error) {
	jsonStr, err := convert.YAMLStringToJSON(yamlStr)
	if err != nil {
		return nil, errorx.NewError(ErrTypeLoaderParseFailed, fmt.Sprintf("YAML 转 JSON 失败: %v", err))
	}

	var messages map[string]map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &messages); err != nil {
		return nil, errorx.NewError(ErrTypeLoaderParseFailed, fmt.Sprintf("JSON 解析失败: %v", err))
	}

	return &YAMLStringLoader{messages: messages}, nil
}

// LoadMessages 加载指定语言的消息
func (y *YAMLStringLoader) LoadMessages(language string) (map[string]string, error) {
	if messages, exists := y.messages[language]; exists {
		return messages, nil
	}
	return nil, errorx.NewError(ErrTypeLoaderLanguageNotFound, fmt.Sprintf("语言: %s", language))
}
