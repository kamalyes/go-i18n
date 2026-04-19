/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-29 00:00:00
 * @FilePath: \go-i18n\errors.go
 * @Description: i18n 错误代码定义
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"github.com/kamalyes/go-toolbox/pkg/errorx"
)

// ==================== i18n 错误类型（6000-6099） ====================

const (
	// ErrTypeLanguageLoadFailed 语言加载失败
	ErrTypeLanguageLoadFailed errorx.ErrorType = 6000 + iota
	// ErrTypeLanguageNotFound 语言未找到
	ErrTypeLanguageNotFound
	// ErrTypeJSONParseFailed JSON 解析失败
	ErrTypeJSONParseFailed
	// ErrTypeMessageLoaderRequired 消息加载器未设置
	ErrTypeMessageLoaderRequired
	// ErrTypeConfigInvalid 配置无效
	ErrTypeConfigInvalid
	// ErrTypeTranslationFailed 翻译失败
	ErrTypeTranslationFailed
)

// init 注册所有错误类型及其默认消息模板
func init() {
	errorx.RegisterError(ErrTypeLanguageLoadFailed, "语言加载失败: %s")
	errorx.RegisterError(ErrTypeLanguageNotFound, "语言未找到: %s")
	errorx.RegisterError(ErrTypeJSONParseFailed, "JSON 解析失败: %s")
	errorx.RegisterError(ErrTypeMessageLoaderRequired, "消息加载器未设置")
	errorx.RegisterError(ErrTypeConfigInvalid, "配置无效: %s")
	errorx.RegisterError(ErrTypeTranslationFailed, "翻译失败: %s")
}
