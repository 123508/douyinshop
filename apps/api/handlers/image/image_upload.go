package image

import (
	"context"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func UploadImage(ctx context.Context, c *app.RequestContext) {
	file, err := c.FormFile("image")

	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"err": err,
		})
	}

	//读取配置文件中的上传类型
	uploadType := config.Conf.AliyunConfig.UploadPath

	var result string

	if uploadType == 1 {
		fileName, err := util.DownloadImages(file)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"err": err,
			})
			return
		}
		result = string(c.Request.Host()) + "/image/get/" + fileName
	} else {
		result, err = util.UploadImagesByIO(file)
	}

	if err != nil {
		messageError := errorno.ParseBasicMessageError(err)

		if messageError.Raw != nil {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"err": err,
			})
			return
		} else {
			c.JSON(messageError.Code, utils.H{
				"err": messageError.Message,
			})
		}

		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"path": result,
	})

}
