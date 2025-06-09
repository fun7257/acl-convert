package acl4ssr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func (ssr *ACL4SSR) convertToClashProxyGroups() string {
	// 使用 strings.Builder 来高效构建字符串
	var builder strings.Builder

	// 首先写入根键和换行符
	builder.WriteString("proxy-groups:\n")
	for _, proxyGroup := range ssr.proxyGroups {
		before, after, found := strings.Cut(proxyGroup, "`")
		if !found {
			continue
		}

		name := before
		builder.WriteString(fmt.Sprintf("  - name: %s\n", name))

		proxyGroup = after
		before, after, found = strings.Cut(proxyGroup, "`")
		if !found {
			continue
		}

		tp := before
		builder.WriteString(fmt.Sprintf("    type: %s\n", tp))

		re := regexp.MustCompile(`^\d+`)

		proxies := []string{}
		includeAll := false
		filter := ""
		url := ""
		interval := 0
		tolerance := 0
		timeout := 0
		for item := range strings.SplitSeq(after, "`") {
			if strings.HasPrefix(item, ".*") {
				includeAll = true
			}

			if strings.HasPrefix(item, "(") {
				includeAll = true
				filter = "(?i)" + strings.TrimPrefix(strings.TrimSuffix(item, ")"), "(")
			}

			if strings.HasPrefix(item, "http") {
				url = item
			}

			if re.MatchString(item) {
				arr := strings.Split(item, ",")
				if len(arr) == 3 {
					interval, _ = strconv.Atoi(arr[0])
					timeout, _ = strconv.Atoi(arr[1])
					tolerance, _ = strconv.Atoi(arr[2])
				}
			}

			if !includeAll && strings.HasPrefix(item, "[]") {
				proxies = append(proxies, strings.TrimPrefix(item, "[]"))
			}
		}

		if includeAll {
			builder.WriteString("    include-all: true\n")
		}

		if filter != "" {
			builder.WriteString(fmt.Sprintf("    filter: %s\n", filter))
		}

		if len(proxies) > 0 {
			builder.WriteString("    proxies:\n")
			for _, proxy := range proxies {
				builder.WriteString(fmt.Sprintf("      - %s\n", proxy))
			}
		}

		if url != "" {
			builder.WriteString(fmt.Sprintf("    url: %s\n", url))
		}

		if interval > 0 {
			builder.WriteString(fmt.Sprintf("    interval: %d\n", interval))
		}

		if tolerance > 0 {
			builder.WriteString(fmt.Sprintf("    tolerance: %d\n", tolerance))
		}

		if timeout > 0 {
			builder.WriteString(fmt.Sprintf("    timeout: %d\n", timeout))
		}

	}

	return builder.String()
}
