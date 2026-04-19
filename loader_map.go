/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-29 00:00:00
 * @FilePath: \go-i18n\loader_map.go
 * @Description: 内存 Map 消息加载器，直接从 map 加载翻译消息
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"fmt"

	gci18n "github.com/kamalyes/go-config/pkg/i18n"
	"github.com/kamalyes/go-toolbox/pkg/errorx"
)

// 确保 MapLoader 同时实现 Loader 和 MessageLoader 接口
var (
	_ Loader               = (*MapLoader)(nil)
	_ gci18n.MessageLoader = (*MapLoader)(nil)
)

// MapLoader 内存 Map 消息加载器
// 直接从内存中的 map 数据结构加载翻译消息
// 适用于程序启动时硬编码或动态构建的翻译数据
type MapLoader struct {
	// messages 语言消息映射
	messages map[string]map[string]string
}

// NewMapLoader 创建内存 Map 消息加载器
// messages: 语言消息映射，格式为 map[语言代码]map[消息键]翻译文本
// 示例: map[string]map[string]string{"zh": {"hello": "你好"}, "en": {"hello": "Hello"}}
func NewMapLoader(messages map[string]map[string]string) *MapLoader {
	return &MapLoader{messages: messages}
}

// LoadMessages 加载指定语言的消息
func (m *MapLoader) LoadMessages(language string) (map[string]string, error) {
	if messages, exists := m.messages[language]; exists {
		return messages, nil
	}
	return nil, errorx.NewError(ErrTypeLoaderLanguageNotFound, fmt.Sprintf("语言: %s", language))
}
