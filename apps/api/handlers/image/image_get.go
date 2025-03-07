package image

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

func GetImage(ctx context.Context, c *app.RequestContext) {
	image := c.Param("image")
	c.File("static/imageStore/" + image)
}
