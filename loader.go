/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-29 00:00:00
 * @FilePath: \go-i18n\loader.go
 * @Description: 消息加载器接口定义与通用工具函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"fmt"

	goi18n "github.com/kamalyes/go-config/pkg/i18n"
	"github.com/kamalyes/go-toolbox/pkg/errorx"
	"github.com/kamalyes/go-toolbox/pkg/mathx"
)

// ==================== 加载器错误类型（6100-6199） ====================

const (
	// ErrTypeLoaderFileNotFound 加载器文件未找到
	ErrTypeLoaderFileNotFound errorx.ErrorType = 6100 + iota
	// ErrTypeLoaderReadFailed 加载器读取失败
	ErrTypeLoaderReadFailed
	// ErrTypeLoaderParseFailed 加载器解析失败
	ErrTypeLoaderParseFailed
	// ErrTypeLoaderLanguageNotFound 加载器语言未找到
	ErrTypeLoaderLanguageNotFound
)

// init 注册加载器错误类型及其默认消息模板
func init() {
	errorx.RegisterError(ErrTypeLoaderFileNotFound, "语言文件未找到: %s")
	errorx.RegisterError(ErrTypeLoaderReadFailed, "读取语言文件失败: %s")
	errorx.RegisterError(ErrTypeLoaderParseFailed, "解析语言数据失败: %s")
	errorx.RegisterError(ErrTypeLoaderLanguageNotFound, "语言未找到: %s")
}

// Loader 消息加载器接口
// 定义从不同数据源加载国际化翻译消息的标准契约
// 与 go-config/pkg/i18n.MessageLoader 接口兼容，所有实现均同时满足两个接口
type Loader interface {
	// LoadMessages 加载指定语言的翻译消息
	// language: 语言代码，如 "zh", "en", "ja" 等
	// 返回该语言的消息映射（键为消息标识，值为翻译文本）
	LoadMessages(language string) (map[string]string, error)
}

// FlattenToMessages 将嵌套 JSON 扁平化为点号格式的消息映射
// 基于 mathx.FlattenMap 实现嵌套结构扁平化，再通过 TransformMapValues 将值转为字符串
// 例如: {"error": {"internal": "错误"}} -> {"error.internal": "错误"}
func FlattenToMessages(data map[string]any) map[string]string {
	ifaceMap := make(map[string]interface{}, len(data))
	for k, v := range data {
		ifaceMap[k] = v
	}

	flat := mathx.FlattenMap(ifaceMap, ".")

	return mathx.TransformMapValues(flat, func(v interface{}) string {
		switch val := v.(type) {
		case string:
			return val
		default:
			return fmt.Sprintf("%v", val)
		}
	})
}

// AsMessageLoader 将 Loader 转换为 go-config 的 MessageLoader 接口
// 由于 Loader 接口与 MessageLoader 接口签名一致，直接类型断言即可
// 如果转换失败返回 nil
func AsMessageLoader(l Loader) goi18n.MessageLoader {
	if ml, ok := l.(goi18n.MessageLoader); ok {
		return ml
	}
	return nil
}
