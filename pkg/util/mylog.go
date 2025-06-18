package util

import log "github.com/sirupsen/logrus"

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "时间",
			log.FieldKeyLevel: "日志类型",
			log.FieldKeyMsg:   "日志内容",
		},
	})
}

// LogError 记录错误日志
func LogError(description, funcName string, err error) {

	param := map[string]interface{}{}
	if description != "" {
		param["问题描述"] = description
	}
	if funcName != "" {
		param["报错函数"] = funcName
	}
	if err != nil {
		param["错误原因"] = err
	}

	field := log.Fields{}
	for k, v := range param {
		field[k] = v
	}

	log.WithFields(field).Error()
}
