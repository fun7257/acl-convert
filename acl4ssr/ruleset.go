package acl4ssr

func (ssr *ACL4SSR) ToClashRule(section, key string) {
	if ssr.Err != nil {
		return
	}

	ssr.loadRuleSet(section, key)
}

func (ssr *ACL4SSR) ToClashRuleProvider(section, key string) {
	if ssr.Err != nil {
		return
	}

	ssr.loadRuleSet(section, key)
}
