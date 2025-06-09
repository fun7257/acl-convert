package acl4ssr

import (
	"net/http"

	"gopkg.in/ini.v1"
)

type ACL4SSR struct {
	INI        *ini.File
	RuleSet    []string
	ProxyGroup []string

	ClashRules         string
	ClashRuleProviders string
	ClashProxyGroups   string

	Err error
}

func (ssr *ACL4SSR) loadRuleSet(section, key string) {
	if ssr.Err != nil {
		return
	}

	ssr.RuleSet = ssr.INI.Section(section).Key(key).ValueWithShadows()
}

func (ssr *ACL4SSR) loadProxyGroup(section, key string) {
	if ssr.Err != nil {
		return
	}

	ssr.ProxyGroup = ssr.INI.Section(section).Key(key).ValueWithShadows()
}

func FetchINI(url string) (ssr *ACL4SSR) {
	resp, err := http.Get(url)
	if err != nil {
		ssr.Err = err
		return
	}
	defer resp.Body.Close()

	ssr.INI, ssr.Err = ini.ShadowLoad(resp.Body)
	return
}
