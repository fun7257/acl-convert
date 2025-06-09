package main

import (
	"fmt"
	"log"

	"github.com/fun7257/acl-convert/acl4ssr"
)

func main() {
	url := "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/refs/heads/master/Clash/config/ACL4SSR_Online_Full_AdblockPlus.ini"
	clash, err := acl4ssr.FetchINI(url).
		LoadProxyGroup("custom", "custom_proxy_group").
		LoadRuleSet("custom", "ruleset").
		Clash()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(clash)
}
