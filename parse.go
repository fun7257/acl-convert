package main

import (
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type RuleSet string

func (rs RuleSet) ConvertToRule() (rule Rule) {
	items := strings.Split(string(rs), ",")
	itemsLen := len(items)
	if itemsLen < 2 { //TODO：定义常量
		return
	}

	rPrefix := "[]"
	rpPrefix := "http"
	proxy := items[0]
	t := items[1]

	if strings.HasPrefix(t, rPrefix) {
		t = strings.TrimPrefix(t, rPrefix)
		switch t {
		case "GEOIP":
			if itemsLen != 3 { //TODO：定义常量
				return
			}
			rule = Rule{
				Type:  "GEOIP",
				Value: items[2],
				Proxy: proxy,
			}
			return
		case "FINAL":
			rule = Rule{
				Type:  "MATCH",
				Proxy: proxy,
			}
			return
		default:
			return
		}
	}

	if strings.HasPrefix(t, rpPrefix) {
		parsedUrl, err := url.Parse(t)
		if err != nil {
			return
		}

		filePath := parsedUrl.Path
		fileNameWithExt := path.Base(filePath)

		// 使用 path.Ext 获取扩展名
		ext := path.Ext(fileNameWithExt)

		// 使用 strings.TrimSuffix 移除扩展名
		fileNameWithoutExt := strings.TrimSuffix(fileNameWithExt, ext)

		rule = Rule{
			Type:  "RULE-SET",
			Value: fileNameWithoutExt,
			Proxy: proxy,
		}

		return
	}

	return
}

func (rs RuleSet) ConvertToRuleProvider() (ruleProvider RuleProvider) {
	items := strings.Split(string(rs), ",")
	itemsLen := len(items)
	if itemsLen < 2 { //TODO：定义常量
		return
	}

	prefix := "http"
	rawUrl := items[1]
	if !strings.HasPrefix(rawUrl, prefix) {
		return
	}

	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return
	}

	filePath := parsedUrl.Path
	fileNameWithExt := path.Base(filePath)

	// 使用 path.Ext 获取扩展名
	ext := path.Ext(fileNameWithExt)

	// 使用 strings.TrimSuffix 移除扩展名
	fileNameWithoutExt := strings.TrimSuffix(fileNameWithExt, ext)

	ruleProvider = RuleProvider{
		Name:     fileNameWithoutExt,
		Type:     "http",
		Url:      rawUrl,
		Proxy:    "DIRECT",
		Behavior: "classical",
		Format:   "text",
		Interval: 24 * 60 * 60,
	}
	return
}

type CustomProxyGroup string

func (cpg CustomProxyGroup) CovertToProxyGroup() (proxyGroup ProxyGroup) {
	items := strings.Split(string(cpg), "`")
	itemsLen := len(items)
	if itemsLen < 2 {
		return
	}

	name := items[0]
	t := items[1]

	includeAll := false
	proxies := []string{}
	use := []string{}
	url := ""
	interval := 0
	tolerance := 0
	timeout := 0
	filter := ""

	proxyPrefix := "[]"

	re := regexp.MustCompile(`^\d+`)

	for _, v := range items[2:] {
		if strings.HasPrefix(v, ".*") {
			includeAll = true
		}

		if strings.HasPrefix(v, "(") {
			includeAll = true
			filter = "(?i)" + strings.TrimPrefix(strings.TrimSuffix(v, ")"), "(")
		}

		if strings.HasPrefix(v, "http") {
			url = v
		}

		if re.MatchString(v) {
			arr := strings.Split(v, ",")
			if len(arr) == 3 {
				interval, _ = strconv.Atoi(arr[0])
				timeout, _ = strconv.Atoi(arr[1])
				tolerance, _ = strconv.Atoi(arr[2])
			}
		}

		if !includeAll && strings.HasPrefix(v, proxyPrefix) {
			proxies = append(proxies, strings.TrimPrefix(v, proxyPrefix))
		}
	}

	proxyGroup.Name = name
	proxyGroup.Type = t
	proxyGroup.IncludeAll = true
	proxyGroup.Url = url
	proxyGroup.Interval = interval
	proxyGroup.Tolerance = tolerance
	proxyGroup.Timeout = timeout
	proxyGroup.Filter = filter

	if len(use) > 0 {
		proxyGroup.Use = use
	}

	if len(proxies) > 0 {
		proxyGroup.Proxies = proxies
	}

	return
}
