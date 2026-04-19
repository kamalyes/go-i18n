/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-29 00:00:00
 * @FilePath: \go-i18n\language.go
 * @Description: 语言解析与标准化工具，集成 go-toolbox/metadata
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"github.com/kamalyes/go-toolbox/pkg/metadata"
)

// NormalizeLanguage 标准化语言代码
// 委托给 metadata.NormalizeLanguage 实现
// 例如: "zh-cn" -> "zh-CN", "zh_CN" -> "zh-CN", "EN" -> "en"
func NormalizeLanguage(lang string) string {
	return metadata.NormalizeLanguage(lang)
}

// ParseAcceptLanguage 解析 Accept-Language 头获取主要语言和地区代码
// 委托给 metadata.ParseAcceptLanguage 实现
// 返回: 语言代码(如 "zh"), 地区代码(如 "CN"), 完整标签(如 "zh-CN")
func ParseAcceptLanguage(acceptLang string) (language, region, fullTag string) {
	return metadata.ParseAcceptLanguage(acceptLang)
}

// ExtractLanguage 从 HTTP 请求中提取语言信息
// 使用 metadata.LanguageExtractor 的默认配置
// 优先级：Query(lang/language) → Header(X-Language) → Cookie(language) → Accept-Language → "en"
func ExtractLanguage(r interface{}) string {
	if req, ok := r.(interface{ GetHTTPRequest() interface{} }); ok {
		httpReq := req.GetHTTPRequest()
		if extractor, ok := httpReq.(interface{ Extract(interface{}) string }); ok {
			return extractor.Extract(httpReq)
		}
	}
	return "en"
}
