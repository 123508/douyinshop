package errorno

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
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

	re0 := regexp.MustCompile(`remote or network error`)
	match := re0.FindStringSubmatch(err.Error())

	if len(match) == 1 {
		fmt.Println(err.Error())
		return &BasicMessageError{Code: 500, Message: "服务出错"}
	}

	// 定义匹配 "code:XXX message:XXX" 的正则表达式
	re := regexp.MustCompile(`code:(\d+)\s+message:(.+)`)
	matches := re.FindStringSubmatch(err.Error())
	if len(matches) != 3 {
		fmt.Println("匹配失败,错误为:", err)
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

func DealWithError(err error, c *app.RequestContext) {

	basicErr := ParseBasicMessageError(err)

	if basicErr.Raw != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": basicErr.Raw.Error(),
		})
	} else {
		c.JSON(basicErr.Code, utils.H{
			"error": basicErr.Message,
		})
	}
}
