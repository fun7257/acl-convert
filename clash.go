package main

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Proxy string `json:"proxy"`
}

type Rules []Rule

func (r Rules) Yaml() (out string) {
	// 使用 strings.Builder 来高效构建字符串
	var builder strings.Builder

	// 首先写入根键和换行符
	builder.WriteString("rules:\n")

	// 遍历所有规则
	for _, v := range r {
		if v.Type == "" || v.Proxy == "" {
			continue
		}

		// 写入每个列表项的缩进和短横线 "- "
		builder.WriteString("  - ")

		// 拼接规则内容
		ruleParts := []string{v.Type}
		if v.Value != "" {
			ruleParts = append(ruleParts, v.Value)
		}
		ruleParts = append(ruleParts, v.Proxy)

		// 将规则的各个部分用逗号连接，并写入 builder
		builder.WriteString(strings.Join(ruleParts, ","))

		// 写入换行符，为下一行做准备
		builder.WriteString("\n")
	}

	// 返回最终构建好的完整字符串
	out = builder.String()
	return

}

type RuleProvider struct {
	// 必须，规则名称，如google,不能重复
	Name string `json:"name" yaml:"-"`
	// 必须，provider类型，可选http / file / inline
	Type string `json:"type" yaml:"type"`
	// 可选，文件路径，不可重复，不填写时会使用 url 的 MD5 作为此文件的文件名
	Path string `json:"path" yaml:"path,omitempty"`
	// 类型为http则必须配置
	Url string `json:"url" yaml:"url,omitempty"`
	// 经过指定代理进行下载/更新
	Proxy string `json:"proxy" yaml:"proxy,omitempty"`
	// 行为，可选domain/ipcidr/classical，对应不同格式的 rule-provider 文件格式，请按实际格式填写
	Behavior string `json:"behavior" yaml:"behavior,omitempty"`
	// 格式，可选 yaml/text/mrs，默认 yaml
	Format string `json:"format" yaml:"format,omitempty"`
	// 更新provider的时间，单位为秒
	Interval int `json:"interval" yaml:"interval,omitempty"`
	// 限制下载文件的最大大小，默认为 0 即不限制文件大小，单位为字节 (b)
	SizeLimit int `json:"sizeLimit" yaml:"sizeLimit,omitempty"`
	// 内容，仅 type 为 inline 时生效
	// TODO: 暂时用不上
	Payload []string `json:"payload" yaml:"payload,omitempty"`
}

type RuleProviders []RuleProvider

func (rps RuleProviders) Yaml() (out string) {
	m := make(map[string]RuleProvider, len(rps))
	for _, v := range rps {
		if v.Name == "" {
			continue
		}

		m[v.Name] = v
	}

	result := struct {
		RuleProviders map[string]RuleProvider `yaml:"rule-providers"`
	}{
		RuleProviders: m,
	}

	buf, err := yaml.Marshal(&result)
	if err != nil {
		return
	}
	out += string(buf)

	return
}

type ProxyGroup struct {
	// 必须，策略组的名字
	Name string `yaml:"name" json:"name"`
	// 必须，策略组的类型
	Type string `yaml:"type" json:"type"`
	// 引入出站代理或其他策略组
	Proxies []string `yaml:"proxies,omitempty" json:"proxies"`
	// 引入代理集合
	Use []string `yaml:"use,omitempty" json:"use"`
	// 健康检查测试地址
	Url string `yaml:"url,omitempty" json:"url"`
	// 健康检查间隔，如不为 0 则启用定时测试，单位为秒
	Interval int `yaml:"interval,omitempty" json:"interval"`
	// 节点切换容差，单位 ms
	Tolerance int `yaml:"tolerance,omitempty" json:"tolerance"`
	// 超时时间，单位ms
	Timeout int `yaml:"timeout,omitempty"`
	// 引入所有出站代理以及代理集合，顺序将按照名称排序
	IncludeAll bool `yaml:"include-all,omitempty" json:"includeAll"`
	// 筛选满足关键词或正则表达式的节点，可以使用 ` 区分多个正则表达式
	Filter string `yaml:"filter,omitempty"`
}

type ProxyGroups []ProxyGroup

func (pgs ProxyGroups) Yaml() (out string) {
	// 使用 strings.Builder 来高效构建字符串
	var builder strings.Builder

	builder.WriteString("proxy-groups:\n")

	for _, v := range pgs {
		if v.Name == "" || v.Type == "" {
			continue
		}

		// 写入每个列表项的缩进和短横线 "- "
		builder.WriteString("  - ")

		builder.WriteString(fmt.Sprintf("name: %s", v.Name))
		builder.WriteString("\n")
		builder.WriteString("    ")

		builder.WriteString(fmt.Sprintf("type: %s", v.Type))
		builder.WriteString("\n")
		builder.WriteString("    ")

		if len(v.Proxies) > 0 {
			var proxyBuilder strings.Builder
			proxyBuilder.WriteString("proxies:\n")
			for _, v := range v.Proxies {
				// 写入每个列表项的缩进和短横线 "- "
				proxyBuilder.WriteString("    ")
				proxyBuilder.WriteString("  - ")
				proxyBuilder.WriteString(v)
				proxyBuilder.WriteString("\n")
			}

			builder.WriteString(proxyBuilder.String())
			builder.WriteString("    ")
		}

		if len(v.Use) > 0 {
			var useBuilder strings.Builder
			useBuilder.WriteString("use:\n")
			for _, v := range v.Proxies {
				// 写入每个列表项的缩进和短横线 "- "
				useBuilder.WriteString("    ")
				useBuilder.WriteString("  - ")
				useBuilder.WriteString(v)
				useBuilder.WriteString("\n")
			}

			builder.WriteString(useBuilder.String())
			builder.WriteString("    ")
		}

		if v.Url != "" {
			builder.WriteString(fmt.Sprintf("url: %s", v.Url))
			builder.WriteString("\n")
			builder.WriteString("    ")
		}

		if v.Interval > 0 {
			builder.WriteString(fmt.Sprintf("interval: %d", v.Interval))
			builder.WriteString("\n")
			builder.WriteString("    ")
		}

		if v.Tolerance > 0 {
			builder.WriteString(fmt.Sprintf("tolerance: %d", v.Tolerance))
			builder.WriteString("\n")
			builder.WriteString("    ")
		}

		if v.Timeout > 0 {
			builder.WriteString(fmt.Sprintf("timeout: %d", v.Timeout))
			builder.WriteString("\n")
			builder.WriteString("    ")
		}

		if v.IncludeAll {
			builder.WriteString("include-all: true")
			builder.WriteString("\n")
			builder.WriteString("    ")
		}

		if v.Filter != "" {
			builder.WriteString(fmt.Sprintf("filter: %s", v.Filter))
			builder.WriteString("\n")
		}

		builder.WriteString("\n")
	}

	out = builder.String()
	return
}
