package test

import (
	"fmt"
	"regexp"
	"testing"
)


func TestMustRegexp(t *testing.T) {
	// 测试正则表达式
	regex := `^\${(.+)}$`
	re := regexp.MustCompile(regex)
	input := "12${prevnode.output.name}"
	matches := re.FindAllStringSubmatch(input, -1)

	if len(matches) == 0 {
		fmt.Println("No matches found.")
		return
	}
	
	for _, match := range matches {
		if len(match) > 1 {
			fmt.Println("Matched group:", match[1])
		}
	}
	
}

func TestRegexp(t *testing.T) {
	// 测试正则表达式
	regex := `^\${(.+)}$`
	re, err := regexp.Compile(regex)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}
	input := "12${prevnode.output.name}"
	matches := re.FindAllStringSubmatch(input, -1)

	if len(matches) == 0 {
		fmt.Println("No matches found.")
		return
	}
	
	for _, match := range matches {
		if len(match) > 1 {
			fmt.Println("Matched group:", match[1])
		}
	}
	
}

func TestFind(t *testing.T) {
	// 测试正则表达式
	regex := `^\${(.+?)}$`
	re := regexp.MustCompile(regex)
	input := "12${prevnode.output.name}000 ${abc}"
	// matches := re.FindAllStringSubmatch(input, -1)
	matches := re.FindStringSubmatch(input)
	fmt.Println("matches:", matches)
	if len(matches) == 0 {
		fmt.Println("No matches found.")
	} else {
		fmt.Println("matched: ", matches[1])
	}

}