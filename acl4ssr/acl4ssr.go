package acl4ssr

import (
	"net/http"
	"net/url"
	"path"
	"strings"

	"gopkg.in/ini.v1"
)

// ACL4SSR 结构体用于处理ACL4SSR配置文件
type ACL4SSR struct {
	iniFile     *ini.File
	rulesets    []string
	proxyGroups []string

	e error
}

// LoadRuleSet 加载规则集
func (ssr *ACL4SSR) LoadRuleSet(section, key string) *ACL4SSR {
	if ssr.e == nil {
		ssr.rulesets = ssr.iniFile.Section(section).Key(key).ValueWithShadows()
	}
	return ssr
}

// LoadProxyGroup 加载代理组
func (ssr *ACL4SSR) LoadProxyGroup(section, key string) *ACL4SSR {
	if ssr.e == nil {
		ssr.proxyGroups = ssr.iniFile.Section(section).Key(key).ValueWithShadows()
	}
	return ssr
}

// Clash 将ACL4SSR配置转换为Clash配置
func (ssr *ACL4SSR) Clash() (out string, err error) {
	if ssr.e != nil {
		err = ssr.e
		return
	}

	out = strings.Join([]string{
		ssr.convertToClashProxyGroups(),
		ssr.convertToClashRuleProviders(),
		ssr.convertToClashRules(),
	}, "\n")
	return
}

// FetchINI 从URL获取INI配置文件
func FetchINI(url string) *ACL4SSR {
	resp, err := http.Get(url)
	if err != nil {
		return &ACL4SSR{
			e: err,
		}
	}
	defer resp.Body.Close()

	ssr := &ACL4SSR{}
	ssr.iniFile, ssr.e = ini.ShadowLoad(resp.Body)
	return ssr
}

// getFileNameFromRawUrl 从URL中获取文件名
func getFileNameFromRawUrl(rawUrl string) (name string) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return
	}

	filePath := parsedUrl.Path
	fileNameWithExt := path.Base(filePath)

	// 使用 path.Ext 获取扩展名
	ext := path.Ext(fileNameWithExt)

	// 使用 strings.TrimSuffix 移除扩展名
	name = strings.TrimSuffix(fileNameWithExt, ext)
	return
}
