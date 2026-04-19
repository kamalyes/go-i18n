/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-29 00:00:00
 * @FilePath: \go-i18n\loader_file.go
 * @Description: 文件系统消息加载器，从磁盘 JSON 文件加载翻译消息
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"fmt"
	"os"
	"path/filepath"

	gci18n "github.com/kamalyes/go-config/pkg/i18n"
	"github.com/kamalyes/go-toolbox/pkg/errorx"
	"github.com/kamalyes/go-toolbox/pkg/json"
)

// 确保 FileLoader 同时实现 Loader 和 MessageLoader 接口
var (
	_ Loader               = (*FileLoader)(nil)
	_ gci18n.MessageLoader = (*FileLoader)(nil)
)

// FileLoader 文件系统消息加载器
// 从指定目录下的 JSON 文件加载翻译消息，支持嵌套 JSON 自动扁平化
// 文件命名规则: {localesPath}/{language}.json，如 ./locales/zh.json
type FileLoader struct {
	// localesPath 翻译文件所在目录路径
	localesPath string
}

// NewFileLoader 创建文件系统消息加载器
// localesPath: 翻译文件所在目录路径，如 "./locales" 或 "resources/locales"
func NewFileLoader(localesPath string) *FileLoader {
	return &FileLoader{
		localesPath: localesPath,
	}
}

// LoadMessages 加载指定语言的消息
// 读取 {localesPath}/{language}.json 文件，支持嵌套 JSON 自动扁平化为点号格式
// language: 语言代码，如 "zh", "en", "ja" 等
func (f *FileLoader) LoadMessages(language string) (map[string]string, error) {
	filePath := filepath.Join(f.localesPath, language+".json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errorx.NewError(ErrTypeLoaderFileNotFound, fmt.Sprintf("路径: %s", filePath))
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errorx.NewError(ErrTypeLoaderReadFailed, fmt.Sprintf("文件: %s, 错误: %v", filePath, err))
	}

	var nested map[string]any
	if err := json.Unmarshal(data, &nested); err != nil {
		return nil, errorx.NewError(ErrTypeLoaderParseFailed, fmt.Sprintf("文件: %s, 错误: %v", filePath, err))
	}

	return FlattenToMessages(nested), nil
}
