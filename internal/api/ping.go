package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func PingHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(200, map[string]string{
		"ping": "pong",
	})
}
