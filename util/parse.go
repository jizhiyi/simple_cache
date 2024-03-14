package util

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseMemorySize 解析内存字符串，返回的是字节数
func ParseMemorySize(sizeStr string) (int64, error) {
	var rate int64 = 1
	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))
	memStrs := []string{"KB", "MB", "GB"}
	for _, memStr := range memStrs {
		if strings.HasSuffix(sizeStr, memStr) {
			rate *= 1024
			sizeStr = strings.TrimSuffix(sizeStr, memStr)
			break
		}
		rate *= 1024
	}
	if rate == 1 {
		return 0, fmt.Errorf("invalid memory size string: %s", sizeStr)
	}
	retSize, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid memory size string: %w", err)
	}
	return retSize * rate, nil
}
