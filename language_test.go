/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-29 00:00:00
 * @FilePath: \go-i18n\language_test.go
 * @Description: 语言解析与标准化工具测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeLanguage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "空字符串", input: "", expected: ""},
		{name: "小写语言", input: "en", expected: "en"},
		{name: "大写语言", input: "EN", expected: "en"},
		{name: "小写地区", input: "zh-cn", expected: "zh-CN"},
		{name: "下划线分隔", input: "zh_CN", expected: "zh-CN"},
		{name: "大写地区", input: "zh-CN", expected: "zh-CN"},
		{name: "带脚本", input: "zh-hans-cn", expected: "zh-Hans-CN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeLanguage(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedLanguage string
		expectedRegion   string
		expectedFullTag  string
	}{
		{name: "空字符串", input: "", expectedLanguage: "", expectedRegion: "", expectedFullTag: ""},
		{name: "简单语言", input: "en", expectedLanguage: "en", expectedRegion: "", expectedFullTag: "en"},
		{name: "语言+地区", input: "zh-CN", expectedLanguage: "zh", expectedRegion: "CN", expectedFullTag: "zh-CN"},
		{name: "带权重", input: "zh-CN;q=0.9,en-US;q=0.8", expectedLanguage: "zh", expectedRegion: "CN", expectedFullTag: "zh-CN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			language, region, fullTag := ParseAcceptLanguage(tt.input)
			assert.Equal(t, tt.expectedLanguage, language)
			assert.Equal(t, tt.expectedRegion, region)
			assert.Equal(t, tt.expectedFullTag, fullTag)
		})
	}
}

func TestExtractLanguage(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{name: "nil", input: nil, expected: "en"},
		{name: "非HTTP请求", input: "string", expected: "en"},
		{name: "无GetHTTPRequest方法", input: struct{}{}, expected: "en"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractLanguage(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
