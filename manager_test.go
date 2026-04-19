/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-29 00:00:00
 * @FilePath: \go-i18n\manager_test.go
 * @Description: 国际化管理器测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"os"
	"testing"

	gci18n "github.com/kamalyes/go-config/pkg/i18n"
	"github.com/kamalyes/go-logger"
	"github.com/stretchr/testify/assert"
)

var (
	testLocalesPath = "locales"
	testFiles       = []string{
		"locales/en.json",
		"locales/zh.json",
	}
)

func TestMain(m *testing.M) {
	setupTestFiles()
	code := m.Run()
	cleanupTestFiles()
	os.Exit(code)
}

func setupTestFiles() {
	os.MkdirAll(testLocalesPath, 0755)
	os.WriteFile(testFiles[0], []byte(`{"hello": "Hello", "greeting": "Hello %s"}`), 0644)
	os.WriteFile(testFiles[1], []byte(`{"hello": "你好", "greeting": "你好 %s"}`), 0644)
}

func cleanupTestFiles() {
	os.RemoveAll(testLocalesPath)
}

func TestNewManager(t *testing.T) {
	tests := []struct {
		name        string
		config      *gci18n.I18N
		expectError bool
	}{
		{name: "使用默认配置", config: nil, expectError: false},
		{name: "使用自定义配置", config: &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}}), EnableFallback: true}, expectError: false},
		{name: "无消息加载器", config: &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en"}, MessageLoader: nil, MessagesPath: ""}, expectError: true},
		{name: "使用 MessagesPath", config: &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en"}, MessageLoader: nil, MessagesPath: testLocalesPath}, expectError: false},
		{name: "带语言映射", config: &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}}), EnableFallback: true, LanguageMapping: map[string]string{"en-US": "en"}}, expectError: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManager(tt.config)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, manager)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, manager)
			}
		})
	}
}

func TestNewManagerWithLogger(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}}), EnableFallback: true}

	tests := []struct {
		name     string
		logger   logger.ILogger
		expected bool
	}{
		{name: "带自定义日志", logger: logger.New(), expected: true},
		{name: "带 nil 日志", logger: nil, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManagerWithLogger(config, tt.logger)
			assert.NoError(t, err)
			assert.NotNil(t, manager)
		})
	}
}

func TestNewManagerWithFormatter(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}}), EnableFallback: true}

	tests := []struct {
		name      string
		formatter Formatter
		expected  bool
	}{
		{name: "带自定义格式化器", formatter: &DefaultFormatter{}, expected: true},
		{name: "带 nil 格式化器", formatter: nil, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManagerWithFormatter(config, tt.formatter)
			assert.NoError(t, err)
			assert.NotNil(t, manager)
		})
	}
}

func TestManager_loadLanguage(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}}), EnableFallback: true}
	manager, err := NewManager(config)
	assert.NoError(t, err)
	assert.NotNil(t, manager)

	err = manager.loadLanguage("en")
	assert.NoError(t, err)
}

func TestManager_GetMessage(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en", "zh"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello", "greeting": "Hello %s"}, "zh": {"hello": "你好", "greeting": "你好 %s"}}), EnableFallback: true}
	manager, err := NewManager(config)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		language string
		key      string
		args     []any
		expected string
	}{
		{name: "英文消息", language: "en", key: "hello", args: nil, expected: "Hello"},
		{name: "中文消息", language: "zh", key: "hello", args: nil, expected: "你好"},
		{name: "带参数的消息", language: "en", key: "greeting", args: []any{"John"}, expected: "Hello John"},
		{name: "不存在的键", language: "en", key: "non_existent", args: nil, expected: "non_existent"},
		{name: "语言回退", language: "fr", key: "hello", args: nil, expected: "Hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.GetMessage(tt.language, tt.key, tt.args...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestManager_GetMessageWithMap(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en", "zh"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"greeting": "Hello {name}"}, "zh": {"greeting": "你好 {name}"}}), EnableFallback: true}
	manager, err := NewManager(config)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		language     string
		key          string
		templateData map[string]any
		expected     string
	}{
		{name: "英文消息带命名参数", language: "en", key: "greeting", templateData: map[string]any{"name": "John"}, expected: "Hello John"},
		{name: "中文消息带命名参数", language: "zh", key: "greeting", templateData: map[string]any{"name": "John"}, expected: "你好 John"},
		{name: "不存在的键", language: "en", key: "non_existent", templateData: map[string]any{"name": "John"}, expected: "non_existent"},
		{name: "语言回退", language: "fr", key: "greeting", templateData: map[string]any{"name": "John"}, expected: "Hello John"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.GetMessageWithMap(tt.language, tt.key, tt.templateData)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestManager_IsLanguageSupported(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en", "zh"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}, "zh": {"hello": "你好"}}), EnableFallback: true}
	manager, err := NewManager(config)
	assert.NoError(t, err)

	tests := []struct {
		language string
		expected bool
	}{
		{language: "en", expected: true},
		{language: "zh", expected: true},
		{language: "fr", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			result := manager.IsLanguageSupported(tt.language)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestManager_GetDefaultLanguage(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en", "zh"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}, "zh": {"hello": "你好"}}), EnableFallback: true}
	manager, err := NewManager(config)
	assert.NoError(t, err)
	assert.Equal(t, "en", manager.GetDefaultLanguage())
}

func TestManager_GetSupportedLanguages(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en", "zh"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}, "zh": {"hello": "你好"}}), EnableFallback: true}
	manager, err := NewManager(config)
	assert.NoError(t, err)
	assert.Len(t, manager.GetSupportedLanguages(), 2)
}

func TestManager_GetConfig(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en", "zh"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}, "zh": {"hello": "你好"}}), EnableFallback: true}
	manager, err := NewManager(config)
	assert.NoError(t, err)
	assert.NotNil(t, manager.GetConfig())
	assert.Equal(t, "en", manager.GetConfig().DefaultLanguage)
}

func TestManager_ReloadLanguage(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}}), EnableFallback: true}
	manager, err := NewManager(config)
	assert.NoError(t, err)
	assert.NoError(t, manager.ReloadLanguage("en"))
}

func TestManager_ResolveLanguage(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en", "zh"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}, "zh": {"hello": "你好"}}), EnableFallback: true}
	manager, err := NewManager(config)
	assert.NoError(t, err)
	assert.Equal(t, "en", manager.ResolveLanguage("en-US"))
}

func TestManager_HasLanguage(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en", "zh"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}, "zh": {"hello": "你好"}}), EnableFallback: true}
	manager, err := NewManager(config)
	assert.NoError(t, err)

	tests := []struct {
		language string
		expected bool
	}{
		{language: "en", expected: true},
		{language: "zh", expected: true},
		{language: "fr", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			assert.Equal(t, tt.expected, manager.HasLanguage(tt.language))
		})
	}
}

func TestManager_GetLoadedLanguages(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en", "zh"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}, "zh": {"hello": "你好"}}), EnableFallback: true}
	manager, err := NewManager(config)
	assert.NoError(t, err)
	assert.Len(t, manager.GetLoadedLanguages(), 2)
}

func TestManager_GetMessageKeys(t *testing.T) {
	config := &gci18n.I18N{Enabled: true, DefaultLanguage: "en", SupportedLanguages: []string{"en", "zh"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello", "greeting": "Hello %s"}, "zh": {"hello": "你好"}}), EnableFallback: true}
	manager, err := NewManager(config)
	assert.NoError(t, err)

	tests := []struct {
		language    string
		expectedLen int
		expectedNil bool
	}{
		{language: "en", expectedLen: 2, expectedNil: false},
		{language: "zh", expectedLen: 1, expectedNil: false},
		{language: "fr", expectedLen: 0, expectedNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			keys := manager.GetMessageKeys(tt.language)
			if tt.expectedNil {
				assert.Nil(t, keys)
			} else {
				assert.Len(t, keys, tt.expectedLen)
			}
		})
	}
}

func TestManager_IsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		enabled  bool
		expected bool
	}{
		{name: "启用国际化", enabled: true, expected: true},
		{name: "禁用国际化", enabled: false, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &gci18n.I18N{Enabled: tt.enabled, DefaultLanguage: "en", SupportedLanguages: []string{"en"}, MessageLoader: NewMapLoader(map[string]map[string]string{"en": {"hello": "Hello"}}), EnableFallback: true}
			manager, err := NewManager(config)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, manager.IsEnabled())
		})
	}
}
