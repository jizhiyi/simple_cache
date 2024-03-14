package main

import (
	"bufio"
	"cache/cache"
	"cache/util"
	"fmt"
	"os"
	"time"
)

type testStruct struct {
	name string
	age  int
}

func main() {
	cache := cache.NewMCache()
	if !cache.SetMaxMemory("10MB") {
		panic("set max memory error")
	}
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	for scanner.Scan() {
		input = scanner.Text()
		args := util.NewArgs(input)
		switch command := args.MustGet(0); command {
		case "set":
			key := args.MustGet(1)
			value := args.MustGet(2)
			cache.Set(key, value, 0*time.Second)
		case "setex":
			key := args.MustGet(1)
			value := args.MustGet(2)
			ex := args.MustInt(3)
			cache.Set(key, value, time.Duration(ex)*time.Second)
		case "setst":
			key := args.MustGet(1)
			cache.Set(key, &testStruct{
				name: "jizhiyi",
				age:  44,
			}, 0*time.Second)
		case "get":
			key := args.MustGet(1)
			if val, ok := cache.Get(key); ok {
				fmt.Printf("%#v\n", val)
			}
		case "del":
			key := args.MustGet(1)
			cache.Del(key)
		case "exists":
			key := args.MustGet(1)
			if cache.Exists(key) {
				fmt.Printf("key: %s exists", key)
			}
		case "flush":
			cache.Flush()
		case "keys":
			fmt.Printf("keys: %d\n", cache.Keys())
		}

	}
}
