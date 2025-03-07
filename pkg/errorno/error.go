package errorno

import (
	"fmt"
	"regexp"
	"strconv"
)

type BasicMessageError struct {
	Code    int
	Message string
	Raw     error
}

func (e *BasicMessageError) Error() string {
	return fmt.Sprintf("code:%d message: %s", e.Code, e.Message)
}

// ParseBasicMessageError 从错误字符串中解析 BasicMessageError
func ParseBasicMessageError(err error) *BasicMessageError {
	// 定义匹配 "code:XXX message:XXX" 的正则表达式
	re := regexp.MustCompile(`code:(\d+)\s+message:(.+)`)
	matches := re.FindStringSubmatch(err.Error())
	if len(matches) != 3 {
		fmt.Println("匹配失败")
		return &BasicMessageError{Raw: err} // 匹配失败
	}

	// 提取 Code 和 Message
	code, _ := strconv.Atoi(matches[1])
	message := matches[2]

	return &BasicMessageError{
		Code:    code,
		Message: message,
	}
}
