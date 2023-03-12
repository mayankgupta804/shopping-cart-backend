package main

import (
	"context"
	"fmt"
	"shopping-cart-backend/config"
	"shopping-cart-backend/internal/api"
	"shopping-cart-backend/internal/domain"
	"shopping-cart-backend/internal/middleware"
	"shopping-cart-backend/internal/repository"
	"shopping-cart-backend/internal/service"
	"shopping-cart-backend/pkg/database"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"
)

func main() {
	config.Load()

	h := server.New(server.WithHostPorts(fmt.Sprintf(":%s", config.App.Server.Port)))

	databaseCfg := database.Config{
		Name:     config.App.Database.Name,
		User:     config.App.Database.User,
		Password: config.App.Database.Password,
		Host:     config.App.Database.Host,
		Port:     config.App.Database.Port,
	}

	db, err := database.NewFromEnv(context.Background(), &databaseCfg)
	if err != nil {
		hlog.Fatal(err)
	}

	var acntRepo domain.AccountRepository
	var acntService service.AccountService

	acntRepo = repository.NewAccountRepository(db)
	acntService = service.NewAccountService(acntRepo)

	var itemRepo domain.ItemRepository
	var itemService service.ItemService

	itemRepo = repository.NewItemRepository(db)
	itemService = service.NewItemService(itemRepo)

	var cartRepo domain.CartRepository
	var cartService service.CartService

	cartRepo = repository.NewCartRepository(db)
	cartService = service.NewCartService(cartRepo, itemRepo)

	adminAuthMiddlware := middleware.NewAuthMiddleware(acntRepo, domain.AdminRole)
	userAuthMiddleware := middleware.NewAuthMiddleware(acntRepo, domain.UserRole)

	regnHandler := api.NewRegistrationHandler(acntService)
	suspendHandler := api.NewSuspendHandler(acntService)
	itemHandler := api.NewItemHandler(itemService)
	cartHandler := api.NewCartHandler(cartService)

	h.GET("/ping", PingHandler)

	if errInit := adminAuthMiddlware.GetInstance().MiddlewareInit(); errInit != nil {
		hlog.Fatalf("authMiddleware.MiddlewareInit() Error: %s", errInit.Error())
	}

	h.NoRoute(adminAuthMiddlware.GetInstance().MiddlewareFunc(), func(ctx context.Context, c *app.RequestContext) {
		claims := jwt.ExtractClaims(ctx, c)
		hlog.Infof("NoRoute claims: %#v\n", claims)
		c.JSON(404, map[string]string{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	account := h.Group("/accounts")
	{
		account.POST("/register", regnHandler.HandleRegistration)
		account.POST("/login", adminAuthMiddlware.GetInstance().LoginHandler)
		account.POST("/logout", adminAuthMiddlware.GetInstance().LogoutHandler)
		account.PUT("/suspend", adminAuthMiddlware.GetInstance().MiddlewareFunc(), suspendHandler.HandleAccountSuspension)
	}

	auth := h.Group("/auth")
	{
		// Refresh time can be longer than token timeout
		auth.GET("/refresh_token", adminAuthMiddlware.GetInstance().RefreshHandler)
	}

	{
		h.GET("/items", userAuthMiddleware.GetInstance().MiddlewareFunc(), itemHandler.HandleGetItem)
		h.POST("/items", adminAuthMiddlware.GetInstance().MiddlewareFunc(), itemHandler.HandleAddItem)
	}

	cart := h.Group("/cart-items")
	{
		cart.POST("/add", userAuthMiddleware.GetInstance().MiddlewareFunc(), cartHandler.HandleAddToCart)
		cart.POST("/remove", userAuthMiddleware.GetInstance().MiddlewareFunc(), cartHandler.HandleRemoveFromCart)

	}

	h.Spin()
}

func PingHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(200, map[string]string{
		"ping": "pong",
	})
}
