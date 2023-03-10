package main

import (
	"context"
	"log"
	_ "shopping-cart-backend/docs"
	"shopping-cart-backend/internal/api"
	"shopping-cart-backend/internal/middleware"
	"shopping-cart-backend/internal/repository"
	"shopping-cart-backend/internal/service"

	"shopping-cart-backend/pkg/database"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
	// "github.com/hertz-contrib/casbin"
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

	cfg := database.Config{
		Name:     "bike_station",
		User:     "mayank",
		Password: "secret",
		Host:     "localhost",
		Port:     "5432",
	}
	db, err := database.NewFromEnv(context.Background(), cfg.DatabaseConfig())
	if err != nil {
		log.Fatal(err)
	}

	acntRepo := repository.NewAccountRepository(db)
	acntService := service.NewAccountService(acntRepo)
	authMiddlware := middleware.NewAuthMiddleware(acntRepo)
	regnHandler := api.NewRegistrationHandler(acntService)
	// casbin.NewCasbinMiddleware()

	url := swagger.URL("http://localhost:8888/swagger/doc.json") // The url pointing to API definition
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))

	// h.GET("/ping", PingHandler)

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	if errInit := authMiddlware.GetInstance().MiddlewareInit(); errInit != nil {
		log.Fatalf("authMiddleware.MiddlewareInit() Error: %s", errInit.Error())
	}

	h.POST("/account", regnHandler.HandleRegistration)
	h.POST("/login", authMiddlware.GetInstance().LoginHandler)
	h.POST("/logout", authMiddlware.GetInstance().LogoutHandler)

	h.NoRoute(authMiddlware.GetInstance().MiddlewareFunc(), func(ctx context.Context, c *app.RequestContext) {
		claims := jwt.ExtractClaims(ctx, c)
		hlog.Infof("NoRoute claims: %#v\n", claims)
		c.JSON(404, map[string]string{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := h.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddlware.GetInstance().RefreshHandler)

	auth.Use(authMiddlware.GetInstance().MiddlewareFunc())
	{
		auth.GET("/ping", PingHandler)
	}
	h.Spin()
}
