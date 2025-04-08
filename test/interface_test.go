package test

import (
	"fmt"
	"math/rand"
	"testing"
)


func TestInterface(t *testing.T) {
	var a interface{}
	random := rand.Intn(10)
	if random > 5 {
		a = "hello"
	} else {
		a = 123
	}
	var b int
	if val, ok := a.(int); ok {
		b = val + 1
		fmt.Println("a is int, incremented b:", b)
	}
	fmt.Println("a", a)
	fmt.Println("b", b)
}

func TestMap(t *testing.T) {
	m := make(map[string]interface{})
	m["key1"] = "value1"
	m["key2"] = 123
	k := "key13"
	fmt.Println(m[k])
}