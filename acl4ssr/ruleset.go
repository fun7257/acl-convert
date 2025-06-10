package acl4ssr

import (
	"fmt"
	"strings"
)

// convertToClashRules 将ACL4SSR规则转换为Clash规则配置
func (ssr *ACL4SSR) convertToClashRules() string {
	// 使用 strings.Builder 来高效构建字符串
	var builder strings.Builder

	// 首先写入根键和换行符
	builder.WriteString("rules:\n")
	for _, ruleset := range ssr.rulesets {
		before, after, found := strings.Cut(ruleset, ",")
		if !found {
			continue
		}

		// 写入每个列表项的缩进和短横线 "- "
		builder.WriteString("  - ")

		if strings.HasPrefix(after, "http") {
			builder.WriteString(strings.Join([]string{
				"RULE-SET",
				getFileNameFromRawUrl(after),
				before,
			}, ",") + "\n")
		}

		if strings.HasPrefix(after, "[]") {
			after = strings.TrimPrefix(after, "[]")
			after = strings.ReplaceAll(after, "FINAL", "MATCH")
			builder.WriteString(strings.Join([]string{
				after,
				before,
			}, ",") + "\n")
		}
	}

	return builder.String()
}

// convertToClashRuleProviders 将ACL4SSR规则转换为Clash规则提供者配置
func (ssr *ACL4SSR) convertToClashRuleProviders() string {
	// 使用 strings.Builder 来高效构建字符串
	var builder strings.Builder

	// 首先写入根键和换行符
	builder.WriteString("rule-providers:\n")
	for _, ruleset := range ssr.rulesets {
		_, after, found := strings.Cut(ruleset, ",")
		if !found {
			continue
		}

		builder.WriteString("  ")
		if strings.HasPrefix(after, "http") {
			builder.WriteString(getFileNameFromRawUrl(after) + ":\n")
			builder.WriteString(fmt.Sprintf("    url: %s\n    behavior: classical\n    interval: 86400\n    format: text\n    type: http\n", after))
		}
	}

	return builder.String()
}
