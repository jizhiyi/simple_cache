package util

import (
	"strconv"
	"strings"
)

type Args struct {
	Args []string
}

func NewArgs(input string) *Args {
	return &Args{Args: strings.Split(input, " ")}
}

func (a *Args) MustGet(index int) string {
	if index < 0 || index >= len(a.Args) {
		return ""
	}
	return a.Args[index]
}

func (a *Args) MustInt(index int) int64 {
	if index < 0 || index >= len(a.Args) {
		return 0
	}
	res, err := strconv.ParseInt(a.Args[index], 10, 64)
	if err != nil {
		res = 0
	}
	return res
}
