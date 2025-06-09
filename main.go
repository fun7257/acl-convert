package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fun7257/acl-convert/acl4ssr"
)

func main() {
	// 1. 定义所有需要的标志
	// -o 用于指定输出文件名
	outputFile := flag.String("o", "", "Output file path. If not specified, prints to stdout.")

	// 自定义更清晰的用法说明
	flag.Usage = func() {
		// 使用 os.Stderr 输出，这是惯例
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <URL>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Fetches the content of a URL and prints it to stdout or a file.\n\n")
		fmt.Fprintf(os.Stderr, "Example: %s -o output.yaml http://example.com\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults() // 打印所有定义的标志及其描述
	}

	// 2. 解析命令行参数
	// flag 包会自动处理 -o 标志，并将剩余的非标志参数收集起来
	flag.Parse()

	// 3. 获取位置参数 (URL)
	// flag.Args() 返回一个包含所有非标志参数的字符串切片
	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "Error: You must specify exactly one URL.")
		flag.Usage() // 打印用法说明
		os.Exit(1)   // 以错误状态退出
	}
	url := flag.Arg(0) // 获取第一个，也是唯一一个非标志参数

	// 4. 拉取acl4ssr配置
	clash, err := acl4ssr.FetchINI(url).
		LoadProxyGroup("custom", "custom_proxy_group").
		LoadRuleSet("custom", "ruleset").
		Clash()

	if err != nil {
		log.Fatal(err)
	}

	// 5. 决定输出目的地 (stdout 或文件)
	var destination io.Writer

	if *outputFile == "" {
		// 如果 -o 未指定，输出到标准输出
		destination = os.Stdout
	} else {
		// 如果 -o 已指定，创建文件
		file, err := os.Create(*outputFile)
		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}
		defer file.Close()
		destination = file
	}

	// 6. 将响应体内容高效地拷贝到目的地
	bytesCopied, err := io.Copy(destination, strings.NewReader(clash))
	if err != nil {
		log.Fatalf("Error writing to destination: %v", err)
	}

	// 如果输出到文件，打印确认信息
	if *outputFile != "" {
		fmt.Fprintf(os.Stderr, "Successfully wrote %d bytes to %s\n", bytesCopied, *outputFile)
	}
}
