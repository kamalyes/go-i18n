/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-29 00:00:00
 * @FilePath: \go-i18n\manager.go
 * @Description: 国际化管理器 - 核心翻译引擎，基于 go-config/i18n 配置
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"fmt"
	"sync"

	gci18n "github.com/kamalyes/go-config/pkg/i18n"
	"github.com/kamalyes/go-logger"
	"github.com/kamalyes/go-toolbox/pkg/errorx"
	"github.com/kamalyes/go-toolbox/pkg/mathx"
)

// 确保 Manager 实现 context.Translator 接口
var _ Translator = (*Manager)(nil)

// Manager 国际化管理器
// 核心翻译引擎，负责消息加载、语言解析和翻译
type Manager struct {
	config    *gci18n.I18N
	messages  map[string]map[string]string
	mutex     sync.RWMutex
	logger    logger.ILogger
	formatter Formatter
}

// NewManager 创建国际化管理器
// config: go-config 的 i18n 配置，如果为 nil 则使用默认配置
func NewManager(config *gci18n.I18N) (*Manager, error) {
	if config == nil {
		config = gci18n.Default()
	}

	if config.MessageLoader == nil && config.MessagesPath != "" {
		config.MessageLoader = NewFileLoader(config.MessagesPath)
	}

	if config.MessageLoader == nil {
		return nil, errorx.NewError(ErrTypeMessageLoaderRequired)
	}

	mgr := &Manager{
		config:    config,
		messages:  make(map[string]map[string]string),
		logger:    logger.New().WithPrefix("[I18N]"),
		formatter: &DefaultFormatter{},
	}

	for _, lang := range config.SupportedLanguages {
		if err := mgr.loadLanguage(lang); err != nil {
			return nil, errorx.NewError(ErrTypeLanguageLoadFailed, fmt.Sprintf("语言 %s: %v", lang, err))
		}
	}

	if config.LanguageMapping != nil {
		for _, targetLang := range config.LanguageMapping {
			if _, exists := mgr.messages[targetLang]; !exists {
				if err := mgr.loadLanguage(targetLang); err != nil {
					mgr.logger.Warn("预加载目标语言 %s 失败: %v", targetLang, err)
					continue
				}
			}
		}
	}

	return mgr, nil
}

// NewManagerWithLogger 创建带自定义日志的国际化管理器
func NewManagerWithLogger(config *gci18n.I18N, log logger.ILogger) (*Manager, error) {
	mgr, err := NewManager(config)
	if err != nil {
		return nil, err
	}
	if log != nil {
		mgr.logger = log
	}
	return mgr, nil
}

// NewManagerWithFormatter 创建带自定义格式化器的国际化管理器
func NewManagerWithFormatter(config *gci18n.I18N, formatter Formatter) (*Manager, error) {
	mgr, err := NewManager(config)
	if err != nil {
		return nil, err
	}
	if formatter != nil {
		mgr.formatter = formatter
	}
	return mgr, nil
}

// loadLanguage 加载指定语言的消息
func (m *Manager) loadLanguage(language string) error {
	messages, err := m.config.MessageLoader.LoadMessages(language)
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.messages[language] = messages

	return nil
}

// GetMessage 获取翻译消息（实现 Translator 接口）
func (m *Manager) GetMessage(language, key string, args ...any) string {
	return m.getMessageInternal(language, key, args, nil)
}

// GetMessageWithMap 使用 map 模板数据获取翻译消息（实现 Translator 接口）
func (m *Manager) GetMessageWithMap(language, key string, templateData map[string]any) string {
	return m.getMessageInternal(language, key, nil, templateData)
}

// getMessageInternal 内部获取消息的实现
func (m *Manager) getMessageInternal(language, key string, args []any, templateData map[string]any) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if message := m.getMessageFromLanguage(language, key); message != "" {
		return m.formatter.Format(message, args, templateData)
	}

	if m.config.EnableFallback && language != m.config.DefaultLanguage {
		if message := m.getMessageFromLanguage(m.config.DefaultLanguage, key); message != "" {
			return m.formatter.Format(message, args, templateData)
		}
	}

	return key
}

// getMessageFromLanguage 从指定语言获取消息
func (m *Manager) getMessageFromLanguage(language, key string) string {
	resolvedLang := m.config.ResolveLanguage(language)
	return m.findMessageInLanguage(resolvedLang, key)
}

// findMessageInLanguage 在指定语言中查找消息
func (m *Manager) findMessageInLanguage(language, key string) string {
	if langMessages, exists := m.messages[language]; exists {
		if message, exists := langMessages[key]; exists {
			return message
		}
	}
	return ""
}

// IsLanguageSupported 检查语言是否被支持（实现 Translator 接口）
func (m *Manager) IsLanguageSupported(language string) bool {
	return m.config.IsSupportedLanguage(language)
}

// GetConfig 获取当前配置
func (m *Manager) GetConfig() *gci18n.I18N {
	return m.config
}

// GetSupportedLanguages 获取支持的语言列表
func (m *Manager) GetSupportedLanguages() []string {
	return m.config.SupportedLanguages
}

// GetDefaultLanguage 获取默认语言
func (m *Manager) GetDefaultLanguage() string {
	return m.config.DefaultLanguage
}

// ReloadLanguage 重新加载指定语言的消息
func (m *Manager) ReloadLanguage(language string) error {
	return m.loadLanguage(language)
}

// ResolveLanguage 解析语言代码（委托给 go-config 的 ResolveLanguage）
func (m *Manager) ResolveLanguage(lang string) string {
	return m.config.ResolveLanguage(lang)
}

// HasLanguage 检查指定语言的消息是否已加载
func (m *Manager) HasLanguage(language string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, exists := m.messages[language]
	return exists
}

// GetLoadedLanguages 获取已加载的语言列表
func (m *Manager) GetLoadedLanguages() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	languages := make([]string, 0, len(m.messages))
	for lang := range m.messages {
		languages = append(languages, lang)
	}
	return languages
}

// GetMessageKeys 获取指定语言的所有消息键
func (m *Manager) GetMessageKeys(language string) []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if langMessages, exists := m.messages[language]; exists {
		keys := make([]string, 0, len(langMessages))
		for key := range langMessages {
			keys = append(keys, key)
		}
		return keys
	}
	return nil
}

// IsEnabled 检查国际化是否启用
func (m *Manager) IsEnabled() bool {
	return mathx.IfDo(m.config != nil, m.config.IsEnabled, false)
}
