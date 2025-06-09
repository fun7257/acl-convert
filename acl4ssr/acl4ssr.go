package acl4ssr

import (
	"net/http"
	"strings"

	"gopkg.in/ini.v1"
)

type ACL4SSR struct {
	iniFile     *ini.File
	rulesets    []string
	proxyGroups []string

	e error
}

func (ssr *ACL4SSR) LoadRuleSet(section, key string) *ACL4SSR {
	if ssr.e == nil {
		ssr.rulesets = ssr.iniFile.Section(section).Key(key).ValueWithShadows()
	}
	return ssr
}

func (ssr *ACL4SSR) LoadProxyGroup(section, key string) *ACL4SSR {
	if ssr.e == nil {
		ssr.proxyGroups = ssr.iniFile.Section(section).Key(key).ValueWithShadows()
	}
	return ssr
}

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
