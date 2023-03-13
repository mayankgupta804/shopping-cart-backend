package main

import (
	"context"
	"fmt"
	"os"
	"shopping-cart-backend/config"
	"shopping-cart-backend/internal/api"
	"shopping-cart-backend/internal/domain"
	"shopping-cart-backend/internal/middleware"
	"shopping-cart-backend/internal/migrations"
	"shopping-cart-backend/internal/repository"
	"shopping-cart-backend/internal/service"
	"shopping-cart-backend/pkg/database"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"
	"github.com/urfave/cli"
)

func main() {
	config.Load()

	clientApp := cli.NewApp()
	clientApp.Name = "Shopping Cart Backend"
	clientApp.Version = "0.0.1"
	clientApp.Commands = []cli.Command{
		{
			Name:        "start:webserver",
			Description: "Start Incident Web Service",
			Action: func(c *cli.Context) {
				StartWebServer()
			},
		},

		{
			Name:        "db:migrate:up",
			Description: "Create migrations",
			Action: func(c *cli.Context) error {
				return migrations.Up(config.App.Database)
			},
		},
		{
			Name:        "db:migrate:down",
			Description: "Destroy migrations",
			Action: func(c *cli.Context) error {
				return migrations.Down(config.App.Database)
			},
		},
	}
	if err := clientApp.Run(os.Args); err != nil {
		panic(err)
	}
}

func StartWebServer() {
	h := server.New(server.WithHostPorts(fmt.Sprintf(":%s", config.App.Server.Port)), server.WithBasePath("/api"))

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
		c.JSON(404, map[string]string{"code": "ROUTE_NOT_FOUND", "message": "Route not found"})
	})

	v1 := h.Group("/v1")

	accounts := v1.Group("/accounts")
	{
		accounts.POST("/register", regnHandler.HandleRegistration)
		accounts.POST("/login", adminAuthMiddlware.GetInstance().LoginHandler)
		accounts.POST("/logout", adminAuthMiddlware.GetInstance().LogoutHandler)
		// CAVEAT: Suspended account will be active until the JWT token expires
		// To resolve the above issue, we need to store the token in Redis OR in process' memory
		// and when a user from a suspended account makes a request, we need to check Redis/memory and block the client
		accounts.PUT("/suspend", adminAuthMiddlware.GetInstance().MiddlewareFunc(), suspendHandler.HandleAccountSuspension)
	}

	auth := v1.Group("/auth")
	{
		// Refresh time can be longer than token timeout
		// TODO: add another middleware to see if the account is not suspended
		// Check DB to see if the account is still active
		auth.GET("/refresh_token", adminAuthMiddlware.GetInstance().RefreshHandler)
	}

	admin := v1.Group("admin")
	{
		admin.POST("/items", adminAuthMiddlware.GetInstance().MiddlewareFunc(), itemHandler.HandleAddItem)
	}

	user := v1.Group("user")
	{
		user.GET("/items", userAuthMiddleware.GetInstance().MiddlewareFunc(), itemHandler.HandleGetItem)
		user.POST("/cart-items", userAuthMiddleware.GetInstance().MiddlewareFunc(), cartHandler.HandleAddToCart)
		user.DELETE("/cart-items", userAuthMiddleware.GetInstance().MiddlewareFunc(), cartHandler.HandleRemoveFromCart)
	}

	h.Spin()
}

func PingHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(200, map[string]string{
		"ping": "pong",
	})
}
