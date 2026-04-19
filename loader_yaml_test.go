/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-29 00:00:00
 * @FilePath: \go-i18n\loader_yaml_test.go
 * @Description: YAML 消息加载器测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testYAMLDir   = "test_yaml"
	testYAMLFiles = []string{
		"test_yaml/en.yaml",
		"test_yaml/zh.yml",
	}
)

func setupYAMLTestFiles() {
	os.MkdirAll(testYAMLDir, 0755)
	os.WriteFile(testYAMLFiles[0], []byte(`hello: Hello
greeting: Hello {name}`), 0644)
	os.WriteFile(testYAMLFiles[1], []byte(`hello: 你好
greeting: 你好 {name}`), 0644)
}

func cleanupYAMLTestFiles() {
	os.RemoveAll(testYAMLDir)
}

func TestNewYAMLLoader(t *testing.T) {
	loader := NewYAMLLoader("./locales")
	assert.NotNil(t, loader)
	assert.Equal(t, "./locales", loader.localesPath)
}

func TestYAMLLoader_LoadMessages(t *testing.T) {
	setupYAMLTestFiles()
	defer cleanupYAMLTestFiles()

	loader := NewYAMLLoader(testYAMLDir)

	tests := []struct {
		name        string
		language    string
		expectError bool
		expectedLen int
	}{
		{name: "加载 .yaml 文件", language: "en", expectError: false, expectedLen: 2},
		{name: "加载 .yml 文件", language: "zh", expectError: false, expectedLen: 2},
		{name: "文件不存在", language: "fr", expectError: true, expectedLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages, err := loader.LoadMessages(tt.language)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, messages)
			} else {
				assert.NoError(t, err)
				assert.Len(t, messages, tt.expectedLen)
			}
		})
	}
}
