package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/ini.v1"
)

func fetchINI(url string) (*ini.File, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ini.ShadowLoad(resp.Body)
}

func main() {
	originName := "oringin.ini"
	_, err := os.Stat(originName)
	if os.IsNotExist(err) {
		cfg, err := fetchINI("https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/refs/heads/master/Clash/config/ACL4SSR_Online_Full_MultiMode.ini")
		if err != nil {
			log.Fatalf("Fail to fetch ini: %v", err)
		}

		err = cfg.SaveTo(originName)
		if err != nil {
			log.Fatal(err)
		}
	}

	cfg, err := ini.ShadowLoad(originName)
	if err != nil {
		log.Fatalf("Fail to load %s: %v", originName, err)
	}

	customSection := cfg.Section("custom")
	if customSection == nil {
		log.Fatal("Custom section not found!")
	}

	// 使用 ValueWithShadows() 方法获取所有同名键的值
	rulesets := customSection.Key("ruleset").ValueWithShadows()
	pgs := customSection.Key("custom_proxy_group").ValueWithShadows()

	ruleProviders := make(RuleProviders, 0, len(rulesets))
	rules := make(Rules, 0, len(rulesets))
	proxyGroups := make(ProxyGroups, 0, len(pgs))
	for _, v := range rulesets {
		ruleProvider := RuleSet(v).ConvertToRuleProvider()
		ruleProviders = append(ruleProviders, ruleProvider)

		rule := RuleSet(v).ConvertToRule()
		rules = append(rules, rule)
	}

	for _, v := range pgs {
		proxyGroup := CustomProxyGroup(v).CovertToProxyGroup()
		proxyGroups = append(proxyGroups, proxyGroup)
	}

	fmt.Println(proxyGroups.Yaml())

	fmt.Println(ruleProviders.Yaml())

	fmt.Println(rules.Yaml())
}
