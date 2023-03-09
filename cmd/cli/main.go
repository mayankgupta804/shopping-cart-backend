package main

import (
	"context"
	_ "shopping-cart-backend/docs"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

// PingHandler
// @summary check to see if the server is running or not
// @Accept application/json
// @Produce application/json
// @Router /ping [get]
func PingHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(200, map[string]string{
		"ping": "pong",
	})
}

// @title Shopping Cart Backend
// @version 1.0
// @description RESTFul HTTP APIs for a hypothetical Shopping Cart.

// @host localhost:8888
// @BasePath /
// @schemes http
func main() {
	h := server.Default()

	h.GET("/ping", PingHandler)

	url := swagger.URL("http://localhost:8888/swagger/doc.json") // The url pointing to API definition
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))

	h.Spin()
}
