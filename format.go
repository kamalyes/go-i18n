/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-29 00:00:00
 * @FilePath: \go-i18n\format.go
 * @Description: 消息格式化工具 - 接口定义与多种格式化策略
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"fmt"
	"strings"
)

// Formatter 消息格式化器接口
// 定义消息格式化的标准契约，支持位置参数和命名模板数据
type Formatter interface {
	// Format 格式化消息
	// message: 消息模板
	// args: 位置参数（如 fmt.Sprintf 风格）
	// templateData: 命名模板数据（如 {key} 风格）
	Format(message string, args []any, templateData map[string]any) string
}

// DefaultFormatter 默认消息格式化器
// 支持三种占位符格式：
//   - fmt.Sprintf 风格（%s, %d 等）→ 位置参数
//   - {key}       → stringx.Format 风格（推荐命名参数）
//   - {{key}}     → 简化模板风格
//   - {{.key}}    → Go template 风格
type DefaultFormatter struct{}

// Format 格式化消息，优先使用 templateData（命名参数），其次使用 args（位置参数）
func (f *DefaultFormatter) Format(message string, args []any, templateData map[string]any) string {
	if templateData != nil {
		return FormatWithTemplateData(message, templateData)
	}

	if len(args) > 0 {
		return fmt.Sprintf(message, args...)
	}

	return message
}

// FormatWithTemplateData 使用模板数据格式化消息
// 支持多种占位符格式，按优先级依次替换
func FormatWithTemplateData(message string, templateData map[string]any) string {
	if templateData == nil {
		return message
	}

	result := message

	for key, value := range templateData {
		valueStr := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, fmt.Sprintf("{{.%s}}", key), valueStr)
		result = strings.ReplaceAll(result, fmt.Sprintf("{{%s}}", key), valueStr)
		result = strings.ReplaceAll(result, fmt.Sprintf("{%s}", key), valueStr)
	}

	result = strings.ReplaceAll(result, ": <no value>", "")
	return result
}

// FormatMessage 便捷函数：使用默认格式化器格式化消息
func FormatMessage(message string, args []any, templateData map[string]any) string {
	return (&DefaultFormatter{}).Format(message, args, templateData)
}
