/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-29 00:00:00
 * @FilePath: \go-i18n\loader_json.go
 * @Description: JSON 消息加载器，从 JSON 字符串加载翻译消息
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"fmt"

	gci18n "github.com/kamalyes/go-config/pkg/i18n"
	"github.com/kamalyes/go-toolbox/pkg/errorx"
	"github.com/kamalyes/go-toolbox/pkg/json"
)

// 确保 JSONLoader 同时实现 Loader 和 MessageLoader 接口
var (
	_ Loader               = (*JSONLoader)(nil)
	_ gci18n.MessageLoader = (*JSONLoader)(nil)
)

// JSONLoader JSON 消息加载器
// 从 JSON 字符串直接解析翻译消息，适用于配置文件内嵌或动态生成的场景
// JSON 格式为: {"zh": {"key": "值"}, "en": {"key": "value"}}
type JSONLoader struct {
	// messages 已解析的语言消息映射
	messages map[string]map[string]string
}

// NewJSONLoader 创建 JSON 消息加载器
// messagesJSON: JSON 格式的翻译消息字符串
// 格式示例: {"zh": {"hello": "你好"}, "en": {"hello": "Hello"}}
func NewJSONLoader(messagesJSON string) (*JSONLoader, error) {
	var messages map[string]map[string]string
	if err := json.Unmarshal([]byte(messagesJSON), &messages); err != nil {
		return nil, errorx.NewError(ErrTypeLoaderParseFailed, fmt.Sprintf("JSON 解析失败: %v", err))
	}

	return &JSONLoader{messages: messages}, nil
}

// LoadMessages 加载指定语言的消息
func (j *JSONLoader) LoadMessages(language string) (map[string]string, error) {
	if messages, exists := j.messages[language]; exists {
		return messages, nil
	}
	return nil, errorx.NewError(ErrTypeLoaderLanguageNotFound, fmt.Sprintf("语言: %s", language))
}
